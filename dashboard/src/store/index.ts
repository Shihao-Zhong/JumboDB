import {combineReducers} from "redux";
import displayReducers from "./displayReducers";
import connectionReducers from "./connectionReducers";
import sampleDataReducers from "./sampleDataReducers";


export default combineReducers({
    display: displayReducers,
    connection: connectionReducers,
    sampleData: sampleDataReducers
});