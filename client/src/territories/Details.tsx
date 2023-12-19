import React from "react";

class LoggingButton extends React.Component {
    handleClick() {
    console.log('test:', this);
}

render() {
    return (
        <button onClick={() => this.handleClick()}>
            Click me
        </button>
        );
    }
}

export default LoggingButton;
