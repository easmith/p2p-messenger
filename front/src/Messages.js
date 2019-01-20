import React, {Component} from 'react';

import Row from "reactstrap/es/Row";
import Col from "reactstrap/es/Col";
import Container from "reactstrap/es/Container";

import PropTypes from 'prop-types';

export default class Messages extends Component {

    scrollToBottom = () => {
        this.messagesEnd.scrollIntoView({ behavior: "smooth" });
    };

    componentDidMount() {
        this.scrollToBottom();
    }

    componentDidUpdate() {
        this.scrollToBottom();
    }

    render() {
        if (!this.props.messages || !this.props.messages.length) {
            return (
                <Container className={"flex-fill"} fluid >
                    <div style={{float: "left", clear: "both"}}
                         ref={(el) => {
                             this.messagesEnd = el;
                         }}>
                    </div>
                </Container>
            )
        }

        return (
            <Container className={"flex-fill"} fluid>
                {this.props.messages.map((m) => {
                    return (
                        <Row className={"border-bottom "}>
                            <Col xs={1}
                                 className={"text-right text-truncate " + (!m.isMine ? "text-danger" : "")}>{m.from}</Col>
                            <Col>{m.content}</Col>
                            <Col xs={1}
                                 className={"text-right " + (!m.isMine ? "text-danger" : "")}>{m.date}</Col>
                        </Row>
                    )
                })}
                <div style={{ float:"left", clear: "both" }}
                     ref={(el) => { this.messagesEnd = el; }}>
                </div>
            </Container>
        )

    }
}


Messages.propTypes = {
    messages: PropTypes.array
};