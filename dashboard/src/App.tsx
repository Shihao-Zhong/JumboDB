import React from 'react';
import './App.css';
import {BrowserRouter as Router, Route, Switch} from "react-router-dom";
import LeftOperationBar from "./components/global/LeftOperationBar";
import Welcome from "./components/pages/Welcome";
import Connection from "./components/pages/Connection";
import SampleData from "./components/pages/SampleData";

export default function App() {

    let header: string = "Welcome";

    return (
        <Router>
            <div>
                <LeftOperationBar header={header}>
                    <Switch>
                        <Route path="/connection">
                            <Connection/>
                        </Route>
                        <Route path="/sampleData">
                            <SampleData/>
                        </Route>
                        <Route path="/">
                            <Welcome/>
                        </Route>
                    </Switch>
                </LeftOperationBar>
            </div>
        </Router>
    );
}

