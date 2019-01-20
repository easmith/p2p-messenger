import React, {Component} from 'react';
import {Navbar, NavbarBrand} from "reactstrap";

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
                <Navbar className={"messenger-navbar"}>
                    <NavbarBrand href="/" className="mr-auto">PeerToPeer Messenger</NavbarBrand>
                </Navbar>
            </header>
        )
    }
}