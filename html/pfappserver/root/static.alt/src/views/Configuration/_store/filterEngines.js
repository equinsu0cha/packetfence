/**
* "$_filter_engines" store module
*/
import Vue from 'vue'
import api from '../_api'

const types = {
  LOADING: 'loading',
  DELETING: 'deleting',
  SUCCESS: 'success',
  ERROR: 'error'
}

// Default values
const state = {
  cache: false, // item details
  message: '',
  itemStatus: ''
}

const getters = {
  isWaiting: state => [types.LOADING, types.DELETING].includes(state.itemStatus),
  isLoading: state => state.itemStatus === types.LOADING
}

const actions = {
  getCollections: ({ state, commit }) => {
    if (state.cache) {
      return Promise.resolve(state.cache).then(collection => Object.values(collection))
    }
    commit('COLLECTIONS_REQUEST')
    return api.filterEnginesCollections().then(response => {
      commit('COLLECTIONS_REPLACED', response.items)
      return Object.values(state.cache)
    }).catch((err) => {
      commit('COLLECTIONS_ERROR', err.response)
      throw err
    })
  },
  getCollection: ({ state, commit, dispatch }, collection) => {
    if (state.cache[collection] && state.cache[collection].items) {
      return Promise.resolve(state.cache[collection]).then(collection => {
        const { items, ...remainder } = collection
        return { ...remainder, ...{ items: Object.values(items) } }
      })
    }
    return dispatch('getCollections').then(() => {
      commit('COLLECTION_REQUEST')
      return api.filterEnginesCollection(collection).then(response => {
        const { items } = response
        commit('COLLECTION_REPLACED', { collection, items })
        const { items: _items, ...remainder } = state.cache[collection]
        return { ...remainder, ...{ items: Object.values(_items) } }
      }).catch((err) => {
        commit('COLLECTION_ERROR', err.response)
        throw err
      })
    })
  },
  getFilterEngine: ({ state, commit, dispatch }, { collection, id }) => {
    if (state.cache[collection] && state.cache[collection][id]) {
      return Promise.resolve(state.cache[collection][id]).then(filterEngine => filterEngine)
    }
    return dispatch('getCollections').then(() => {
      commit('ITEM_REQUEST')
      const { [collection]: { resource } = {} } = state.cache
      return api.filterEngine({ resource, id }).then(item => {
        commit('ITEM_REPLACED', { collection, id, item })
        return state.cache[collection][id]
      }).catch((err) => {
        commit('ITEM_ERROR', err.response)
        throw err
      })
    })
  },
  options: ({ commit, dispatch }, { collection, id }) => {
    if (id) {
      return dispatch('getCollections').then(() => {
        commit('ITEM_REQUEST')
        const { [collection]: { resource } = {} } = state.cache
        return api.filterEngineOptions({ resource, id }).then(response => {
          commit('ITEM_SUCCESS')
          return response
        }).catch((err) => {
          commit('ITEM_ERROR', err.response)
          throw err
        })
      })
    } else {
      commit('ITEM_REQUEST')
      return api.filterEnginesOptions(collection).then(response => {
        commit('ITEM_SUCCESS')
        return response
      }).catch((err) => {
        commit('ITEM_ERROR', err.response)
        throw err
      })
    }
  },
  createFilterEngine: ({ commit }, { collection, data }) => {
    commit('ITEM_REQUEST')
    return api.createFilterEngine({ collection, data }).then(response => {

console.log('createFilterEngine', { collection, data, response })

      commit('ITEM_SUCCESS')
      return response
    }).catch((err) => {
      commit('ITEM_ERROR', err.response)
      throw err
    })
  },
  updateFilterEngine: ({ commit, dispatch }, { collection, id, data }) => {
    return dispatch('getCollections').then(() => {
      commit('ITEM_REQUEST')
      const { [collection]: { resource } = {} } = state.cache
      return api.updateFilterEngine({ resource, id, data }).then(response => {

console.log('updateFilterEngine', { collection, id, data, response })

        commit('ITEM_SUCCESS')
        return response
      }).catch((err) => {
        commit('ITEM_ERROR', err.response)
        throw err
      })
    })
  },
  deleteFilterEngine: ({ commit, dispatch }, { collection, id }) => {
    return dispatch('getCollections').then(() => {
      commit('ITEM_REQUEST')
      const { [collection]: { resource } = {} } = state.cache
      return api.deleteFilterEngine({ resource, id }).then(response => {

console.log('deleteFilterEngine', { collection, id, response })

        commit('ITEM_SUCCESS')
        return response
      }).catch((err) => {
        commit('ITEM_ERROR', err.response)
        throw err
      })
    })
  },
  sortFilterEngines: ({ commit }, { collection, data }) => {
    const params = {
      items: data
    }
    commit('ITEM_REQUEST', types.LOADING)
    return api.sortFilterEngines({ collection, params }).then(response => {
      commit('ITEM_SUCCESS')
      return response
    }).catch(err => {
      commit('ITEM_ERROR', err.response)
      throw err
    })
  }
}

const mutations = {
  COLLECTIONS_REQUEST: (state, type) => {
    state.itemStatus = type || types.LOADING
    state.message = ''
  },
  COLLECTIONS_REPLACED: (state, items) => {
    state.itemStatus = types.SUCCESS
    Vue.set(state, 'cache', items.reduce((items, item) => {
      const { collection } = item
      items[collection] = item
      return items
    }, {}))
  },
  COLLECTIONS_ERROR: (state, response) => {
    state.itemStatus = types.ERROR
    if (response && response.data) {
      state.message = response.data.message
    }
  },

  COLLECTION_REQUEST: (state, type) => {
    state.itemStatus = type || types.LOADING
    state.message = ''
  },
  COLLECTION_REPLACED: (state, { collection, items }) => {
    state.itemStatus = types.SUCCESS
    if (!(collection in state.cache)) {
      Vue.set(state.cache, collection, {})
    }
    Vue.set(state.cache[collection], 'items', items.reduce((items, item) => {
      const { id } = item
      items[id] = item
      return items
    }, {}))
  },
  COLLECTION_ERROR: (state, response) => {
    state.itemStatus = types.ERROR
    if (response && response.data) {
      state.message = response.data.message
    }
  },

  ITEM_REQUEST: (state, type) => {
    state.itemStatus = type || types.LOADING
    state.message = ''
  },
  ITEM_REPLACED: (state, { collection, id, item }) => {
    state.itemStatus = types.SUCCESS
    if (!(collection in state.cache)) {
      Vue.set(state.cache, collection, {})
    }
    Vue.set(state.cache[collection], id, item)
  },
  ITEM_SUCCESS: (state) => {
    state.itemStatus = types.SUCCESS
    state.message = ''
  },
  ITEM_ERROR: (state, response) => {
    state.itemStatus = types.ERROR
    if (response && response.data) {
      state.message = response.data.message
    }
  }
}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
}