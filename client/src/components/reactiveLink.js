import React, { PureComponent } from 'react'
import { PropTypes } from 'prop-types'
import { NavLink } from 'react-router-dom'

export default class ReactiveLink extends PureComponent {
	static propTypes = {
		disabled: PropTypes.bool.isRequired,
		// text: PropTypes.string.isRequired,
		children: PropTypes.oneOfType([PropTypes.arrayOf(PropTypes.node), PropTypes.node]).isRequired,
	}

	render() {
		const { disabled, children, ...rest } = this.props

		// return disabled ? <span>{text}</span> : <NavLink {...rest}>{text}</NavLink>
		return disabled ? <span>{children}</span> : <NavLink {...rest}>{children}</NavLink>
	}
}
