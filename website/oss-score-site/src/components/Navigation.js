import React from 'react'
import { Navbar, Container, Nav} from 'react-bootstrap'
import './Navigation.css';

import './Navigation.css';

const Navigation = () => {
    return (
        <>
            <Navbar className='navigation-bar' expand="lg" variant="dark">
                <Container>
<<<<<<< HEAD
                <Navbar.Brand href="/">OSS-Score</Navbar.Brand>
                    <Nav className="me-auto ml-0">
=======
                    <Nav className="m-auto">
>>>>>>> 64195d8ab0b4843997990ef1d35834df716f8346
                        <Nav.Link class="nav-home" href="/">Home</Nav.Link>
                        <Nav.Link class="nav-link" href="/about">About</Nav.Link>
                        <Nav.Link class="nav-link" href="/extension">Extension</Nav.Link>
                        <Nav.Link class="nav-link" href="/generate-scores">How We Generate Score</Nav.Link>
                        <Nav.Link class="nav-link" href="/privacypolicy">Privacy Policy</Nav.Link>
                    </Nav>
                </Container>
            </Navbar>
        </>
    );
}

export default Navigation