import { combineReducers } from 'redux'
import { routeReducer } from 'redux-simple-router'
import unbalance from './unbalance'

export default combineReducers({
  unbalance,
  router: routeReducer
})
