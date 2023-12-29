import Container from 'react-bootstrap/Container';
import Nav from 'react-bootstrap/Nav';
import Navbar from 'react-bootstrap/Navbar';
import geopolitical from "../assets/geopolitics.png"

function HeaderComponent() {
    return (
        <Navbar expand="lg" bg="dark" data-bs-theme="dark" sticky="top">
            <Container>
                <Navbar.Brand href="/" className="me-5">
                    <img
                        alt="icon geopolitical simulation"
                        src={geopolitical}
                        width="30"
                        height="30"
                        className="d-inline-block align-top me-1"
                    />{' '}
                    Geopolitical Simulation
                </Navbar.Brand>
                <Navbar.Toggle aria-controls="basic-navbar-nav" />
                <Navbar.Collapse id="basic-navbar-nav">
                    <Nav className="me-auto">
                        <Nav.Link href="/">Territories</Nav.Link>
                        <Nav.Link href="/countries">Countries</Nav.Link>
                        <Nav.Link href="/market">Market</Nav.Link>
                    </Nav>
                </Navbar.Collapse>
            </Container>
        </Navbar>
    );
}

export default HeaderComponent;