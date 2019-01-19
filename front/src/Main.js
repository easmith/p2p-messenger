import React, {Component} from 'react';
import {Button, Col, Container, Row} from "reactstrap";

import Peers from "./Peers";
import MessageInput from "./MessageInput";
import Messages from "./Messages";


import update from 'immutability-helper';

export default class Main extends Component {


    constructor(props, context) {
        super(props, context);

        this.state = {
            socket: null,
            iam: null,
            interlocutor: null,
            peers: {},
            messages: {},
        }
    }

    componentDidMount() {
        let socket = new WebSocket("ws://" + document.location.hostname + (document.location.port ? ':' + document.location.port : '') + "/ws");

        this.setState({socket: socket}, () => {
            socket.onopen = function () {
                console.log("Соединение установлено.");
                socket.send(JSON.stringify({cmd: "HELLO"}));
                socket.send(JSON.stringify({cmd: "PEERS"}));
            };

            socket.onmessage = this.onMessage;

            socket.onclose = function (event) {
                if (event.wasClean) {
                    console.log('Соединение закрыто чисто');
                } else {
                    console.log('Обрыв соединения');
                }
                console.log('Код: ' + event.code + ' причина: ' + event.reason);
            };

            socket.onerror = function (error) {
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

        switch (parsedMessage.cmd) {
            case "NAME" : {
                this.setState({iam: {name: parsedMessage.name, id: parsedMessage.id}})
                break;
            }
            case "PEERS" : {
                let peers = {};
                parsedMessage.peers.forEach((p) => {
                    let v = this.state.peers[p.id];
                    p.counter = v ? v.counter : 0;
                    peers[p.id] = p
                });
                this.setState({peers: peers});
                return;
            }
            case "MESS" : {
                let peerId = "";
                let fromName = "";
                let counter = 1;

                if (parsedMessage.from === this.state.iam.id) {
                    // это наше сообщение
                    peerId = parsedMessage.to;
                    fromName = this.state.iam.name;
                } else {
                    // это сообщение от другого пира
                    peerId = parsedMessage.from;
                    let peer = this.state.peers[peerId];
                    if (peer) {
                        fromName = peer.name;
                        counter = peer.counter + 1;
                    } else {
                        fromName = parsedMessage.from.substr(0, 10);
                    }
                }

                let oldMessages = this.state.messages[peerId];
                if (!oldMessages) {
                    oldMessages = []
                }

                let message = {
                    date: new Date().toLocaleTimeString(['ru-RU', 'en-US'], {hour12: false}),
                    isMine: parsedMessage.from === this.state.iam.id,
                    from: fromName,
                    content: parsedMessage.content
                };

                oldMessages.push(message);

                console.log(oldMessages)

                this.setState({
                    peers: update(this.state.peers, {[peerId]: {counter: {$set: counter}}}),
                    messages: update(this.state.messages, {[peerId]: {$set: oldMessages}})
                });
                break;
            }
            default : {
                console.warn("Unknown cmd: " + parsedMessage.cmd)
            }
        }

    };

    updatePeers = () => {
        this.state.socket.send(JSON.stringify({cmd: "PEERS"}));
    };

    sendMessage = (msg) => {
        let cmd = JSON.stringify({
            cmd: "MESS",
            from: this.state.iam.id,
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
                        <Messages
                            messages={this.state.interlocutor ? this.state.messages[this.state.interlocutor.id] : []}/>
                        <MessageInput interlocutor={this.state.interlocutor} onSendMessage={this.sendMessage}/>
                    </Col>
                </Row>
            </Container>
        )
    }
}