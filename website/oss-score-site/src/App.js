import './App.css';
import {Route, Routes} from "react-router-dom";
import React  from 'react';
import Home from './components/Homepage.js';
import About from './components/About.js';
import Navigation from './components/Navigation.js';
import Extension from './components/Extension.js';
import GenerateScores from './components/GenerateScores.js';
import PrivacyPolicy from './components/PrivacyPolicy.js';
import NotFound from './components/NotFound.js';
import 'bootstrap/dist/css/bootstrap.min.css';

function App() {
  return (
    <div className="nap">
      {/* Maybe put nav in separate file: https://medium.com/swlh/responsive-navbar-using-react-bootstrap-5e0e0bd33bd6 */}
      < Navigation />
      < Routes >
          <Route exact path='/' element={<Home/>} />
          <Route exact path='/about' element= {<About/>}/>
          <Route exact path='/extension' element={<Extension/>}/>
          <Route exact path='/generate-scores' element={<GenerateScores/>}/>
          <Route exact path='/privacypolicy' element={<PrivacyPolicy/>}/>

          {/* Page Not Found Routes */}
          <Route path="" element={<NotFound />} />
          <Route path="*" element={<NotFound />} />
          <Route element={<NotFound />} />

        </Routes>
      </div>  
  );
}

export default App;
