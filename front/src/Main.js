import React, {Component} from 'react';
import {Button, Col, Container, Row} from "reactstrap";

import Peers from "./Peers";
import MessageInput from "./MessageInput";
import Messages from "./Messages";

export default class Main extends Component {


    constructor(props, context) {
        super(props, context);

        this.state = {
            socket: null,
            interlocutor: null,
            peers: [],
            messages: [],
        }
    }

    componentDidMount () {
        let socket = new WebSocket("ws://"+ document.location.hostname+(document.location.port ? ':'+document.location.port : '') + "/ws");

        this.setState({socket: socket}, () => {
            socket.onopen = function() {
                console.log("Соединение установлено.");
                socket.send(JSON.stringify({cmd:"PEERS"}));
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
        });
    }

    onMessage = (event) => {
        console.log("Получены данные " + event.data);
        let parsedMessage = JSON.parse(event.data);

        if (!parsedMessage.cmd) {
            console.error("something wrong with data");
            return;
        }

        if (parsedMessage.cmd === "PEERS") {
            this.setState({peers: parsedMessage.peers})
        }

        if (parsedMessage.cmd === "MESS") {
            this.setState(
                {
                    messages: [...this.state.messages, parsedMessage.content]
                }
            )
            //{message: parsedMessage.content})
        }


    };

    updatePeers = () => {
        this.state.socket.send(JSON.stringify({cmd:"PEERS"}));
    };

    sendMessage = (msg) => {
        let cmd = JSON.stringify({
            cmd: "MESS",
            to: this.state.interlocutor.id,
            content: msg
        });
        this.state.socket.send(cmd);
    };

    selectPeer = (peer) => {
        this.setState({
            interlocutor: peer
        })
    };


    render() {
        let interlocutorName = this.state.interlocutor ? " with " + this.state.interlocutor.name : "";
        return (
            <Container className={"vh-100 mt-4"} fluid>
                <Row>
                    <Col className={"border-right"}>
                        <h3>Peers <Button color="info" size={"sm"} onClick={this.updatePeers}>update</Button></h3>
                    </Col>
                    <Col xs={9}>
                        <h3>Chat {interlocutorName}</h3>
                    </Col>
                </Row>
                <Row className={"h-75"}>
                    <Col xs={3} className={"border-right"}>
                        <Peers peers={this.state.peers} onSelectPeer={this.selectPeer}/>
                    </Col>
                    <Col xs={9}>
                        <Messages messages={this.state.messages}/>
                        <MessageInput interlocutor={this.state.interlocutor} onSendMessage={this.sendMessage}/>
                    </Col>
                </Row>
            </Container>
        )
    }
}