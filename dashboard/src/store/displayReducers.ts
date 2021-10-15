import * as actionTypes from "./actionTypes";


const initialState: any = {
    header: "Welcome"
};

const reducer = (
    state: any = initialState,
    action: any
) => {
    switch (action.type) {
        case actionTypes.UPDATE_UI_HEADER: {
            return {
                ...state,
                header: action.header
            };
        }
        default:
            return state;
    }
};

export default reducer;


