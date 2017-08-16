import { createStore, combineReducers } from 'redux';

const authReducer = (state={}, action) => {
  switch (action.type) {
    case 'AUTH':
      state = { ...state, auth: 1 };
      break;
    case 'UNAUTH':
      state = { ...state, auth: 0 };
      break;
  }
  return state;
};

const userReducer = (state=null, action) => {
  switch (action.type) {
    case 'USER':
      state = { ...state, user: action.payload };
      break;
    default:
      state = { ...state, user: 'anonymous' };
  }
  return state;
};

const reducers = combineReducers({
  auth: authReducer,
  user: userReducer,
});

export const store = createStore(reducers);

export const auth = () => new Promise((res, rej) => {
    store.subscribe(() => {
      const state = store.getState();
      res(state);
    });
  });

// store.subscribe(() => {
//   console.log(store.getState());
// });

// store.dispatch({ type: 'AUTH' });
// store.dispatch({ type: 'USER', payload: 'Ekrem' });
