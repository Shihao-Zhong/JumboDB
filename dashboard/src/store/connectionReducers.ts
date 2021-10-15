import * as actionTypes from "./actionTypes";


const initialState: any = {
    currentIndex: -1,
    allHosts: []
};

const reducer = (
    state: any = initialState,
    action: any
) => {
    switch (action.type) {
        case actionTypes.UPDATE_HOST: {
            return {
                ...state,
                currentIndex: action.payload
            };
        }

        case actionTypes.ADD_HOST: {
            if (state.allHosts.includes(action.payload)) {
                return {
                    ...state,
                    currentHost: action.payload
                };
            } else {
                return {
                    ...state,
                    allHosts: [...state.allHosts, action.payload],
                    currentIndex: action.payload
                };
            }
        }

        default:
            return state;
    }
};

export default reducer;


