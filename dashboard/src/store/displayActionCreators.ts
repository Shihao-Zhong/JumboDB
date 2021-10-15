import * as actionTypes from "./actionTypes";

export function updateHeader(header: string) {
    const action = {
        type: actionTypes.UPDATE_UI_HEADER,
        header: header,
    };
    return (dispatch: any) => {
        dispatch(action);
    }
}