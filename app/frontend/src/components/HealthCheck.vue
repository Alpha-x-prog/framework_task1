<script setup lang="ts">
import { ref, onMounted } from 'vue'
const status = ref<'idle'|'ok'|'fail'>('idle')
const body = ref('')

onMounted(async () => {
  try {
    const res = await fetch('/api/healthz')
    body.value = await res.text()
    status.value = res.ok ? 'ok' : 'fail'
  } catch {
    status.value = 'fail'
  }
})
</script>

<template>
  <section>
    <h2>Health check</h2>
    <p>Status: <strong>{{ status }}</strong></p>
    <pre style="white-space:pre-wrap">{{ body }}</pre>
  </section>
</template>
