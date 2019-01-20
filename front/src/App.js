import React, {Component} from 'react';
import './App.scss';
import Header from "./Header";
import Main from "./Main";

class App extends Component {

    render() {
        return (
            <div className={"d-flex flex-column h-100"}>
                <Header/>
                <Main/>
            </div>
        );
    }

}

export default App;
