<!--

<pf-sidebar-item :item="section" :filter="filter" indent>

-->
<template>
  <div>
    <b-nav-item
      exact-active-class="active"
      v-if="visible"
      v-bind="$attrs"
      :to="item.path"
      :key="item.name"
      >
      <div class="pf-sidebar-item" :class="{ 'ml-3': indent }">
        <div>
          <text-highlight :queries="[filter]">{{ item.name }}</text-highlight>
          <text-highlight class="figure-caption text-nowrap" v-if="item.caption" :queries="[filter]">{{ item.caption }}</text-highlight>
        </div>
        <icon class="mx-1" :name="item.icon" v-if="item.icon"></icon>
        <slot/>
      </div>
    </b-nav-item>
    <b-nav class="pf-sidenav mb-2" v-if="showSavedSearches && savedSearches.length > 0" vertical>
      <div class="pf-sidenav-group" v-t="'Saved Searches'"></div>
      <b-nav-item
        exact-active-class="active"
        v-for="search in savedSearches"
        :key="search.name"
        :to="search.route"
        class="saved-search"
      >
        <div class="pf-sidebar-item ml-3">
          <text-highlight :queries="[filter]">{{ search.name }}</text-highlight>
          <icon class="mx-1" name="trash-alt" role="button" @click.stop.prevent="deleteSavedSearch(search)"></icon>
        </div>
      </b-nav-item>
    </b-nav>
  </div>
</template>

<script>
import TextHighlight from 'vue-text-highlight'

export default {
  name: 'pf-sidebar-item',
  components: {
    TextHighlight
  },
  props: {
    item: {
      default: { name: 'undefined', path: '/' }
    },
    filter: {
      default: ''
    },
    indent: {
      type: Boolean,
      default: false
    }
  },
  data () {
    return {
      visible: {
        type: Boolean,
        default: true
      }
    }
  },
  computed: {
    showSavedSearches () {
      const { item: { saveSearchNamespace } = {} } = this
      return this.visible && saveSearchNamespace
    },
    savedSearches () {
      return this.$store.getters['saveSearch/cache'][this.item.saveSearchNamespace] || []
    }
  },
  methods: {
    deleteSavedSearch (search) {
      const { item: { saveSearchNamespace } = {} } = this
      this.$store.dispatch('saveSearch/remove', { namespace: saveSearchNamespace, search: { name: search.name } })
    }
  },
  mounted () {
    if ('can' in this.item) {
      this.visible = this.$can.apply(null, this.item.can.split(' '))
    }
    /*
    if ('saveSearchNamespace' in this.item) {
      this.$store.dispatch('saveSearch/get', this.item.saveSearchNamespace).then(savedSearches => {
        this.$set(this, 'savedSearches', savedSearches)
      })
    }
    */
  }
}
</script>

<style lang="scss">
@import '../styles/variables';

.saved-search {
  a {
    svg.fa-icon {
      visibility: hidden;
    }
  }
}
.saved-search:hover a,
.saved-search a.active {
  svg.fa-icon {
    visibility: visible;
    color: rgba($body-bg, .7);
    &:hover {
      color: $body-bg;
    }
  }
}
</style>
