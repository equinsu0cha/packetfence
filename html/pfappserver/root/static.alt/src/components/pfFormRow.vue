<template>
    <b-row :class="rowClass" :align-v="alignV">
        <b-col sm="3" class="col-form-label"><label class="mb-0" :class="labelClass" :for="id" v-text="columnLabel"></label></b-col>
        <b-col sm="9"><slot/></b-col>
    </b-row>
</template>

<script>
export default {
  name: 'pf-form-row',
  props: {
    id: {
      type: String
    },
    rowClass: {
      type: String,
      default: 'form-row'
    },
    columnLabel: {
      type: String
    },
    labelClass: {
      type: String
    },
    alignV: {
      type: String,
      default: 'center'
    }
  },
  mounted () {
    // Propagate id to the first element of the default slot
    if (this.id && this.$slots.default().length) {
      this.$slots.default()[0].elm.id = this.id
    }
  }
}
</script>
