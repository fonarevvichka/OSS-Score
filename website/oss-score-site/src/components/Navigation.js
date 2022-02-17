import React from 'react'
import { Navbar, Container, Nav} from 'react-bootstrap'

import './Navigation.css';

const Navigation = () => {
    return (
        <>
            <Navbar expand="lg" bg="primary" variant="dark">
                <Container>
                    <Nav className="ml-auto">
                        <Nav.Link class="nav-home" href="/">Home</Nav.Link>
                        <Nav.Link class="nav-link" href="/about">About</Nav.Link>
                        <Nav.Link class="nav-link" href="/extension">Extension</Nav.Link>
                        <Nav.Link class="nav-link" href="/generate-scores">How We Generate Score</Nav.Link>
                        <Nav.Link class="nav-link" href="/categories">Categories</Nav.Link>
                        <Nav.Link class="nav-link" href="/workforus">Work For Us</Nav.Link>
                    </Nav>
                </Container>
            </Navbar>
        </>
    );
}

export default Navigation