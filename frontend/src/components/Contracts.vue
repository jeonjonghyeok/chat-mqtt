<template>
  <div class="contract">
    <button class="contract-id">
      {{contract.id}}
      {{contract.sign_on}}
      {{contract.sender_id}}
    </button>
  </div>
</template>

<script>
export default {
  name: "Contracts",
  props: {
    contract: Object
  }
};
</script>

<style scoped>
.contract{
  font-family: sans-serif;
  margin: 5px 0;
}
.contract-id {
  font-size: 25px;
}
</style>