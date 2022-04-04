//import logo from './logo.svg';
import './App.css';
import {Route, Routes} from "react-router-dom";
import React  from 'react';

// import {Navbar, Container, Nav, NavDropdown} from 'react-bootstrap'

import Home from './components/Homepage.js';
import About from './components/About.js';
import Navigation from './components/Navigation.js';
import Extension from './components/Extension.js';
import GenerateScores from './components/GenerateScores.js';
import PrivacyPolicy from './components/PrivacyPolicy.js';
import NotFound from './components/NotFound.js';
import backgroundSVG from './images/subtle-prism.svg';


import 'bootstrap/dist/css/bootstrap.min.css';
//import FixedNavbarExample from './navbar.js'

function App() {
  return (
<<<<<<< HEAD
    // style={{ backgroundImage: `url(${backgroundSVG}`}}
    <div className="app" >
      <div className="nap" style={{}}>
        {/* Maybe put nav in separate file: https://medium.com/swlh/responsive-navbar-using-react-bootstrap-5e0e0bd33bd6 */}
        < Navigation />
        < Routes >
          <Route exact path='/' element={<Home />} />
          <Route exact path='/about' element={<About />} />
          <Route exact path='/extension' element={<Extension />} />
          <Route exact path='/generate-scores' element={<GenerateScores />} />
          <Route exact path='/categories' element={<Categories />} />
          <Route exact path='/workforus' element={<WorkForUs />} />
          <Route exact path='/privacypolicy' element={<PrivacyPolicy />} />

=======
    <div className="nap">
      {/* Maybe put nav in separate file: https://medium.com/swlh/responsive-navbar-using-react-bootstrap-5e0e0bd33bd6 */}
      < Navigation />
      < Routes >
          <Route exact path='/' element={<Home/>} />
          <Route exact path='/about' element= {<About/>}/>
          <Route exact path='/extension' element={<Extension/>}/>
          <Route exact path='/generate-scores' element={<GenerateScores/>}/>
          <Route exact path='/privacypolicy' element={<PrivacyPolicy/>}/>
>>>>>>> 64195d8ab0b4843997990ef1d35834df716f8346

          {/* Page Not Found Routes */}
          <Route path="" element={<NotFound />} />
          <Route path="*" element={<NotFound />} />
          <Route element={<NotFound />} />

        </Routes>
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

        {/* < svg view-box="0 0 1600 900" >
    <path fill="#5b4693" opacity="1" d="M0,297C267,403,534,368,801,577,C1068,786,1335,391,1602,451,C1600, 900,1600, 900,1600, 900C1600, 900,1600, 900,1600, 900C1600, 900,1600, 900,1600, 900C1600, 900,1600, 900,1600, 900L1600,900C1333,900,1066,900,799,900,C532,900,265,900,-2,900,C0, 900,0, 900,0, 900C0, 900,0, 900,0, 900C0, 900,0, 900,0, 900C0, 900,0, 900,0, 900L1401,900L0,900Z" />
    </svg > */}
      </div>
    </div>
  
  );
}

export default App;
