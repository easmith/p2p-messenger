import React, {Component} from 'react';

import Row from "reactstrap/es/Row";
import Col from "reactstrap/es/Col";
import Container from "reactstrap/es/Container";

import PropTypes from 'prop-types';

export default class Messages extends Component {

    render() {
        if (this.props.messages) {
            return (
                <Container className={"messages rounded"} fluid>
                    {this.props.messages.map((m) => {
                        return (
                            <Row>
                                <Col xs={1}
                                     className={"text-right border-bottom " + (!m.isMine ? "text-danger" : "")}>{m.date}</Col>
                                <Col xs={1}
                                     className={"text-right border-bottom " + (!m.isMine ? "text-danger" : "")}>{m.from}</Col>
                                <Col>{m.content}</Col>
                            </Row>
                        )
                    })}
                </Container>
            )
        } else {
            return (
                <Container className={"messages rounded"} fluid>
                </Container>
            )
        }
    }
}


Messages.propTypes = {
    messages: PropTypes.object
};