import * as actionTypes from "./actionTypes";


const initialState: any = {
    allData: []
};

const reducer = (
    state: any = initialState,
    action: any
) => {
    switch (action.type) {
        case actionTypes.UPDATE_ALL_DATA: {
            return {
                ...state,
                allData: action.payload
            };
        }
        default:
            return state;
    }
};

export default reducer;


