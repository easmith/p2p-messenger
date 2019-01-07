import React, {Component} from 'react';

import PropTypes from 'prop-types';
import Row from "reactstrap/es/Row";
import Col from "reactstrap/es/Col";
import Container from "reactstrap/es/Container";

export default class Messages extends Component {

    render() {
        return (
            <Container className={"messages rounded"} fluid>
                {this.props.messages.map((m) => {
                    return (
                        <Row className={"danger"}>
                            <Col xs={1} className={"text-right border-bottom text-danger"}>from</Col>
                            <Col> {m}</Col>
                        </Row>
                    )
                })}

            </Container>
        )
    }
}



Messages.propTypes = {
    messages: PropTypes.array
};