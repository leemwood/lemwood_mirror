package browser

import (
    "errors"
    "fmt"
    "net/url"
    "regexp"
    "strings"

    "github.com/gocolly/colly/v2"
)

// ResolveRepoURL 访问给定的源 URL 并尝试查找 GitHub 仓库链接。
// 如果源已经是 github.com URL，则直接返回。
// repoSelector: "regex:<pattern>" 将使用正则表达式匹配锚点 href。
// 否则视为 CSS 选择器，使用第一个匹配元素的 href。
func ResolveRepoURL(source string, repoSelector string) (string, error) {
    if source == "" {
        return "", errors.New("源 url 为空")
    }
    u, err := url.Parse(source)
    if err != nil {
        return "", fmt.Errorf("无效的源 url: %w", err)
    }

    // 如果源看起来像直接的 GitHub 仓库 URL (https://github.com/owner/repo)，则按原样返回
    if strings.Contains(u.Host, "github.com") {
        parts := strings.Split(strings.Trim(u.Path, "/"), "/")
        if len(parts) >= 2 && parts[0] != "" && parts[1] != "" {
            // Direct repo URL
            return "https://github.com/" + parts[0] + "/" + parts[1], nil
        }
        // 否则，它是 GitHub 页面（例如搜索/结果）。我们将爬取下面的锚点。
    }

    c := colly.NewCollector(
        colly.MaxDepth(1),
        colly.AllowedDomains(u.Host),
    )
    var found string
    var re *regexp.Regexp
    cssSelector := "a"
    if repoSelector != "" {
        if strings.HasPrefix(repoSelector, "regex:") {
            pattern := strings.TrimPrefix(repoSelector, "regex:")
            compiled, err := regexp.Compile(pattern)
            if err != nil {
                return "", fmt.Errorf("repo_selector 中的正则表达式无效: %w", err)
            }
            re = compiled
            cssSelector = "a"
        } else {
            cssSelector = repoSelector
        }
    }

    // 默认严格匹配：https://github.com/<owner>/<repo> 或 /<owner>/<repo>
    defaultAbsRe := regexp.MustCompile(`^https://github\.com/[^/]+/[^/#?]+$`)
    defaultRelRe := regexp.MustCompile(`^/[^/]+/[^/#?]+$`)

    c.OnHTML(cssSelector, func(e *colly.HTMLElement) {
        if found != "" {
            return
        }
        href := strings.TrimSpace(e.Attr("href"))
        if href == "" {
            return
        }
        // 如果提供了自定义正则表达式
        if re != nil {
            if re.MatchString(href) {
                if strings.HasPrefix(href, "/") {
                    found = "https://github.com" + href
                } else {
                    found = href
                }
            }
            return
        }
        // 否则对 GitHub 仓库 URL 应用默认匹配
        if defaultAbsRe.MatchString(href) {
            found = href
            return
        }
        if defaultRelRe.MatchString(href) {
            found = "https://github.com" + href
            return
        }
    })
    if err := c.Visit(source); err != nil {
        return "", fmt.Errorf("访问源失败: %w", err)
    }
    if found == "" {
        return "", errors.New("未从源页面找到 github 仓库 url")
    }
    return found, nil
}
