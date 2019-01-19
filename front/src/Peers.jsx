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
                {
                    Object.keys(this.props.peers).map((id) => {
                        return (
                            <ListGroupItem color={"secondary"} key={id}>
                                <Button color={"default"} onClick={this.selectItem}
                                        data-name={this.props.peers[id].name}
                                        data-id={id}>{this.props.peers[id].name}</Button>
                                <Badge
                                    className={"float-right"}>{this.props.peers[id].counter > 0 ? this.props.peers[id].counter : ""}</Badge>
                            </ListGroupItem>
                        )
                    })
                }
            </ListGroup>
        )
    }
}

Peers.propTypes = {
    peers: PropTypes.object,
    onSelectPeer: PropTypes.func
};