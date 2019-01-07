import React, {Component} from 'react';
import {Collapse, Nav, Navbar, NavbarBrand, NavbarToggler, NavItem, NavLink} from "reactstrap";

export default class Header extends Component {

    constructor(props) {
        super(props);

        this.toggleNavbar = this.toggleNavbar.bind(this);
        this.state = {
            collapsed: true
        };
    }

    toggleNavbar() {
        this.setState({
            collapsed: !this.state.collapsed
        });
    }

    render() {
        return (
            <header>
                <Navbar color={"dark"} dark >
                    <NavbarBrand href="/" className="mr-auto">PeerToPeer Messenger</NavbarBrand>
                    <NavbarToggler onClick={this.toggleNavbar} className="mr-2" />
                    <Collapse isOpen={!this.state.collapsed} navbar>
                        <Nav navbar>
                            <NavItem>
                                <NavLink href="https://github.com/easmith/p2p-messenger">Source</NavLink>
                            </NavItem>
                            <NavItem>
                                @Copyleft easmith
                            </NavItem>
                        </Nav>
                    </Collapse>
                </Navbar>
            </header>
        )
    }
}