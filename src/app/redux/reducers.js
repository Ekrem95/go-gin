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
      return state = 0;
  }
};

export const store = createStore(reducer, 6);
