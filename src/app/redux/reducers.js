import { createStore } from 'redux';

const reducer = (state, action) => {
  switch (action.type) {
    case 'AUTH':
      return state = 1;
      break;
    case 'UNAUTH':
      return state = 0;
      break;
    default:
      return state;
  }
};

export const store = createStore(reducer, 6);

export const auth = () => new Promise((res, rej) => {
    store.subscribe(() => {
      const state = store.getState();
      res(state);
    });
  });
