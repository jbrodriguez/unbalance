import React from 'react'
import { connect } from 'react-redux'
import { Link } from 'react-router'
import { actions } from '../redux/modules/unbalance'
import styles from './HomeView.scss'

// We define mapStateToProps where we'd normally use
// the @connect decorator so the data requirements are clear upfront, but then
// export the decorated component after the main class definition so
// the component can be tested w/ and w/o being connected.
// See: http://rackt.github.io/redux/docs/recipes/WritingTests.html
const mapStateToProps = (state) => ({
	unraid: state.unraid
})
export class HomeView extends React.Component {
	static propTypes = {
		unraid: React.PropTypes.object.isRequired,
		calculate: React.PropTypes.func.isRequired,
		move: React.PropTypes.func.isRequired
	}

	render () {
		return (
		<div className='container text-center'>
		<h1>Welcome to the React Redux Starter Kit</h1>
		<h2>
		Sample Counter:&nbsp;
		<span className={styles['counter--green']}>{this.props.counter}</span>
		</h2>
		<button className='btn btn-default'
		onClick={() => this.props.increment(1)}>
		Increment
		</button>
		<button className='btn btn-default'
		onClick={this.props.doubleAsync}>
		Double (Async)
		</button>
		<hr />
		<Link to='/about'>Go To About View</Link>
		</div>
		)
	}
}

export default connect(mapStateToProps, actions)(HomeView)
