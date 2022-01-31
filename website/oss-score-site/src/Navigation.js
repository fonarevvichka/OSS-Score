import React from 'react'
import { Navbar, Container, Nav, NavDropdown } from 'react-bootstrap'

const Navigation = () => {
    return (
        <>
            <Navbar bg="primary" variant="dark">
                <Container>
                    <Nav className="me-auto">
                        <Nav.Link href="/">Home</Nav.Link>
                        <Nav.Link href="/about">About</Nav.Link>
                        <Nav.Link href="/extension">Extension</Nav.Link>
                        <Nav.Link href="/generate-scores">How We Generate Score</Nav.Link>
                        <Nav.Link href="/accomplishments">Accomplishments</Nav.Link>
                        <Nav.Link href="/work-with-us">Work With Us</Nav.Link>
                    </Nav>
                </Container>
            </Navbar>
        </>
    );
}

export default Navigation