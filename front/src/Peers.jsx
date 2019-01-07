import React, {Component} from 'react';
import {Badge, Button, Col, ListGroup, ListGroupItem} from "reactstrap";
import PropTypes from 'prop-types';


export default class Peers extends Component {
    static defaultProps = {
        peers: []
    };

    render() {
        return (
            <Col xs={3}>
                <ListGroup className={"peers"}>
                    {this.props.peers.map(p => {
                        return (
                            <ListGroupItem key={p.id}>
                                <Button color={"default"}>{p.name}</Button>
                                <Badge className={"float-right"}></Badge>
                            </ListGroupItem>
                        )
                    })}
                </ListGroup>
            </Col>

        )
    }
}

Peers.propTypes = {
    peers: PropTypes.array,
    onSelectPeer: PropTypes.func
};