import React, { Component, PropTypes } from 'react'

export default class Provider extends Component {
	constructor(props, context) {
		super(props, context)

		console.log('provider.props: ', props)
		console.log('provider.context: ', context)

		this.model = props.model
		this.dispatch = props.dispatch
	}

	getChildContext() {
		return {
			model: this.model,
			dispatch: this.dispatch,
		}
	}

	render() {
		return this.props.children
	}
}
Provider.propTypes = {
	model: PropTypes.object.isRequired,
	dispatch: PropTypes.func.isRequired,
	children: PropTypes.element.isRequired,
}
Provider.childContextTypes = {
	model: PropTypes.object.isRequired,
	dispatch: PropTypes.func.isRequired,
}