import React from 'react'
import { Navbar, Container, Nav} from 'react-bootstrap'
import './Navigation.css';

import './Navigation.css';
import smallLogo from '../images/favicon.ico'

const Navigation = () => {
    return (
            <Navbar className='navigation-bar' expand="md" variant="dark">
                <Navbar.Brand href="/" style={{paddingLeft:"10px"}}>
                    <img className='navLogo' src={smallLogo}></img>
                    OSS-Score</Navbar.Brand>
                    <Nav className="navbar-nav ms-auto">
                        <Nav.Link class="nav-home" href="/">Home</Nav.Link>
                        <Nav.Link class="nav-link" href="/about">About</Nav.Link>
                        <Nav.Link class="nav-link" href="/extension">Extension</Nav.Link>
                        <Nav.Link class="nav-link" href="/generate-scores">How We Generate Score</Nav.Link>
                        <Nav.Link class="nav-link" href="/privacypolicy">Privacy Policy</Nav.Link>
                    </Nav>
            </Navbar>
    );
}

export default Navigation