import React, {Component} from 'react';
import {Button, Col, Container, Input, InputGroup, InputGroupAddon, Row} from "reactstrap";

import Peers from "./Peers";

export default class Main extends Component {


    constructor(props, context) {
        super(props, context);

        this.state = {
            socket: null,
            messages: ""
        }
    }

    componentDidMount ()  {
        let socket = new WebSocket("ws://localhost:35035/ws");

        this.setState({socket: socket});

        socket.onopen = function() {
            console.log("Соединение установлено.");
            socket.send(JSON.stringify({cmd: "HELLO"}))
        };

        socket.onmessage = this.onMessage;

        socket.onclose = function(event) {
            if (event.wasClean) {
                console.log('Соединение закрыто чисто');
            } else {
                console.log('Обрыв соединения');
            }
            console.log('Код: ' + event.code + ' причина: ' + event.reason);
        };

        socket.onerror = function(error) {
            console.log("Ошибка " + error.message);
        };

    }

    onMessage = (event) => {
        let parsedMessage = JSON.parse(event.data);
        this.setState({
            messages:  this.state.messages + "\n" + parsedMessage.cmd
        });

        console.log("Получены данные " + event.data);
    }


    _handleEnter = (e) => {
        if (e.key === 'Enter') {
            this.state.socket.send(e.target.value);
            e.target.value = "";
        }
    };

    updatePeers = () => {
        this.state.socket.send(JSON.stringify({cmd:"PEERS"}));
    };


    render() {
        return (
            <Container className={"vh-100 mt-3"}>
                <Row>
                    <Col>
                        <h3>Peers <Button color="info" size={"sm"} onClick={this.updatePeers}>update</Button></h3>
                    </Col>
                    <Col xs={9}>
                        <h3>Chat</h3>
                    </Col>
                </Row>
                <Row className={"h-75"}>
                    <Peers peers={[{name:"name", id: "id"}]}/>
                    <Col xs={9}>
                        <pre style={{border: "1px solid red"}}>{this.state.messages}</pre>
                    </Col>
                </Row>
                <Row className={"mt-3"}>
                    <Col>
                        <InputGroup>
                            <InputGroupAddon addonType="prepend">
                                <Button color="danger">Select peer</Button>
                            </InputGroupAddon>
                            <Input placeholder="Type a message and press Enter" onKeyPress={this._handleEnter} />
                            <InputGroupAddon addonType="append">
                                <Button color="success">Send!</Button>
                            </InputGroupAddon>
                        </InputGroup>
                    </Col>
                </Row>
            </Container>
        )
    }
}