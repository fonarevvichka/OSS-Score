import logo from './logo.svg';
import './App.css';
import {Route, Routes} from "react-router-dom";
import React, { Component }  from 'react';

// import {Navbar, Container, Nav, NavDropdown} from 'react-bootstrap'

import Home from './homepage.js';
import About from './about.js';
import Navigation from './Navigation.js';


import 'bootstrap/dist/css/bootstrap.min.css';
//import FixedNavbarExample from './navbar.js'

function App() {
  return (
    <div className="nap">
      {/* Maybe put nav in separate file: https://medium.com/swlh/responsive-navbar-using-react-bootstrap-5e0e0bd33bd6 */}
      < Navigation />
      < Routes >
        <Route path='/' component={Home}/>
        <Route path='/about' component={About}/>
        {/*<Route exact path='/extension' component={Extension}/>
        <Route exact path='/generate-scores' component={GenerateScores}/>
        <Route exact path='/accomplishments' component={Accomplishments}/>
  <Route exact path='/work-for-us' component={WorkForUs}/>*/}
      </Routes>
      < Home />
      {/* <Route exact path="/homepage.js" component={Home} />  */}
      {/* <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <p>
          Edit <code>src/App.js</code> and save to reload.
        </p>
        <a
          className="App-link"
          href="https://reactjs.org"
          target="_blank"
          rel="noopener noreferrer"
        >
          Learn React
        </a>
      </header>*/}

      
  </div>
  );
}

export default App;
