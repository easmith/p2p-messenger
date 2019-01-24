import React, {Component} from 'react';
import {Badge, ListGroup, ListGroupItem} from "reactstrap";
import PropTypes from 'prop-types';


export default class Peers extends Component {
    static defaultProps = {
        peers: []
    };

    selectItem = (id, key) => {
        return (e) => {
            this.props.onSelectPeer({
                id: id,
                name: key
            });
        }

        // console.log(elem);
        // console.log(elem.target);
        // console.log(elem.parentNode);
        // console.log(elem.target.getAttribute("data-id"));
        // this.props.onSelectPeer({
        //     id: elem.getAttribute("data-id"),
        //     name: elem.getAttribute("data-name")
        // });
    };

    render() {
        if (!Object.keys(this.props.peers).length) {
            return (
                <div>empty :(</div>
            )
        }
        return (
            <ListGroup flush className={"peers mt-3"}>
                {
                    Object.keys(this.props.peers).map((id) => {
                        return (
                            <ListGroupItem className={"d-flex align-content-center flex-wrap"}
                                           key={id}
                                           data-name={this.props.peers[id].name}
                                           data-id={id}
                                           onClick={this.selectItem(id, this.props.peers[id].name)}>
                                <div className={"p-2 pl-3"}>
                                    {this.props.peers[id].name}
                                </div>
                                <div className={"p-2 pr-3 ml-auto"}>
                                    <Badge>{this.props.peers[id].counter > 0 ? this.props.peers[id].counter : ""}</Badge>
                                </div>
                            </ListGroupItem>
                        )
                    })
                }
            </ListGroup>
        )
    }
}

Peers.propTypes = {
    selectedPeerId: PropTypes.object,
    peers: PropTypes.object,
    onSelectPeer: PropTypes.func
};