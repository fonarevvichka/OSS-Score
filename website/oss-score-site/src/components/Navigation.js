import React from 'react'
import { Navbar, Container, Nav} from 'react-bootstrap'
import './Navigation.css';

import './Navigation.css';
import smallLogo from '../images/favicon.ico'

const Navigation = () => {
    return (
        // <div className='container'>
        <nav className="navbar navbar-expand-sm navbar-dark" style={{backgroundColor:'#5b4da2'}}>
                <a href='/# ' className="navbar-brand mb-0 h1">
                        <img className='d-inline-block align-top' src={smallLogo} width="30" height="30" alt='OSS-Logo'></img>
                        OSS-Score
                </a>
                <button type="button" data-bs-toggle="collapse" data-bs-target="#navbarNav" className="navbar-toggler" aria-controls="navbarNav" aria-expanded="false" aria-label="Toggle navigation">
                    <span class="navbar-toggler-icon"></span>
                </button>
            <div className="collapse navbar-collapse" style={{display: 'inline-block'}} id="navbarNav">
                    <ul classsName="nav navbar-nav">
                        <li class="nav-item active">
                            <a href="/# " className="nav-link">Home</a>
                        </li>
                        <li class="nav-item active">
                                <a href="/about " className="nav-link">About</a>
                        </li>
                        <li class="nav-item active">
                                <a href="/extension " className="nav-link">Extension</a>
                        </li>
                        <li class="nav-item active">
                                <a href="/generate-scores " className="nav-link">Score Generation</a>
                        </li>
                        <li class="nav-item active">
                                <a href="/privacypolicy " className="nav-link">Privacy Policy</a>
                        </li>
                    </ul>
                </div>
            </nav>
        // </div>


        // <>
        //     <Navbar className='navigation-bar' expand="lg" variant="dark">
        //         <Container>
        //         <Navbar.Brand href="/">
        //             <img className='navLogo' src={smallLogo}></img>
        //             OSS-Score</Navbar.Brand>
        //             {/* <Nav className="me-auto ml-0"> */}
        //             <Nav className="navbar-nav ms-auto">
        //                 {/* <Navbar.Brand href="/">
        //                     <img className='navLogo' src={smallLogo}></img>
        //                     OSS-Score</Navbar.Brand> */}
        //                 {/* <Nav.Link class="nav-home" href="/">Home</Nav.Link> */}
        //                 <Nav.Link class="nav-link" href="/about">About</Nav.Link>
        //                 <Nav.Link class="nav-link" href="/extension">Extension</Nav.Link>
        //                 <Nav.Link class="nav-link" href="/generate-scores">How We Generate Score</Nav.Link>
        //                 <Nav.Link class="nav-link" href="/privacypolicy">Privacy Policy</Nav.Link>
        //             {/* </Nav>    */}
        //             </Nav>
        //         </Container>
        //     </Navbar>
        // </>
    );
}

export default Navigation