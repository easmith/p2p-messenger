import React, {Component} from 'react';
import {Button, Input, InputGroup, InputGroupAddon} from "reactstrap";

import PropTypes from 'prop-types';

export default class MessageInput extends Component {

    _handleEnter = (e) => {
        if (e.key === 'Enter') {
            this.props.onSendMessage(e.target.value);
            e.target.value = "";
        }
    };

    render() {
        if (this.props.interlocutor == null) {
            return <div className={"mt-3"}>Select peer</div>
        }
        return (
            <InputGroup className={"mt-3"}>
                <InputGroupAddon addonType="prepend">
                    <Button color="danger">{this.props.interlocutor.name}</Button>
                </InputGroupAddon>
                <Input placeholder="Type a message and press Enter" onKeyPress={this._handleEnter}/>
                {/*<InputGroupAddon addonType="append">*/}
                {/*<Button color="success">Send!</Button>*/}
                {/*</InputGroupAddon>*/}
            </InputGroup>
        )
    }
}


MessageInput.propTypes = {
    interlocutor: PropTypes.object,
    onSendMessage: PropTypes.func,
};