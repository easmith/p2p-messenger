import React, {Component} from 'react';

import PropTypes from 'prop-types';

export default class MessageInput extends Component {
    _handleEnter = (e) => {
        if (e.key === 'Enter') {
            this.props.onSendMessage(e.target.value);
            e.target.value = "";
        }
    };

    onSend = () => {
        this.props.onSendMessage(this.inputTag.value);
        this.inputTag.value = "";
    };

    render() {
        if (this.props.interlocutor == null) {
            return (
                <div className={"message-input"}>
                    <div>Select a peer</div>
                </div>
            )
        }

        return (
            <div className={"message-input d-flex"}>
                <input key={"msgInput"} placeholder="Type a message and press Enter" onKeyPress={this._handleEnter}
                       ref={(el) => { this.inputTag = el; }}/>
                <button onClick={this.onSend}>Send</button>
            </div>
        )
    }
}


MessageInput.propTypes = {
    interlocutor: PropTypes.object,
    onSendMessage: PropTypes.func,
};