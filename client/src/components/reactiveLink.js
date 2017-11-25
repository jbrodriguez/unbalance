import React, { PureComponent } from 'react'
import { PropTypes } from 'prop-types'
import { NavLink } from 'react-router-dom'

export default class ReactiveLink extends PureComponent {
	static propTypes = {
		disabled: PropTypes.bool.isRequired,
		text: PropTypes.string.isRequired,
	}

	render() {
		const { disabled, text, ...rest } = this.props

		return disabled ? <span>{text}</span> : <NavLink {...rest}>{text}</NavLink>
	}
}
