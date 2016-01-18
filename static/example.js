var MessageElement = React.createClass({
  render: function() {
    return (
      <li>{this.props.message}</li>
    )
  }
});

var Display = React.createClass({


  getInitialState: function() {
    return {messages: []};
  },
  componentDidMount: function() {
    this.setState({messages: ["okay"]});
  },
  render: function() {
    var rendered_messages = this.state.messages.map(function(message, id) {
      return (
        <MessageElement message={message} key={id} />
      )
    });
    return (
      <div>
        <p>Hi!</p>
        <div>
          <ol>{rendered_messages}</ol>
        </div>
      </div>
    );
  }
});

ReactDOM.render(
  <Display url="/api/bands" />,
  document.getElementById('content')
);
