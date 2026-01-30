<script setup>
import DefaultLayout from '@/layouts/DefaultLayout.vue'
import CookiesConsent from '@/components/CookiesConsent.vue'
</script>

<template>
  <DefaultLayout>
    <router-view />
  </DefaultLayout>
  <CookiesConsent />
</template>
