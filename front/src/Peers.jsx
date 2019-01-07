import React, {Component} from 'react';
import {Badge, Button, ListGroup, ListGroupItem} from "reactstrap";
import PropTypes from 'prop-types';


export default class Peers extends Component {
    static defaultProps = {
        peers: []
    };

    selectItem = (peer) => {
        console.log(peer.target.getAttribute("data-id"));

        this.props.onSelectPeer({
            id: peer.target.getAttribute("data-id"),
            name: peer.target.getAttribute("data-name")
        });
    };

    render() {
        return (
            <ListGroup className={"peers"}>
                {this.props.peers.map(p => {
                    return (
                        <ListGroupItem color={"secondary"} key={p.id}>
                            <Button color={"default"} onClick={this.selectItem} data-name={p.name} data-id={p.id}>{p.name}</Button>
                            <Badge className={"float-right"}></Badge>
                        </ListGroupItem>
                    )
                })}
            </ListGroup>
        )
    }
}

Peers.propTypes = {
    peers: PropTypes.array,
    onSelectPeer: PropTypes.func
};