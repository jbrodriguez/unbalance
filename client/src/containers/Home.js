import React, { Component, PropTypes } from 'react'
import { connect } from 'react-redux'
import { Link } from 'react-router'
import { actions as counterActions } from '../redux/modules/counter'
import styles from './HomeView.scss'

// example state
// state = {
// 	config: {
// 		folders: [
// 			"movies/films",
// 			"movies/tvshows"
// 		],
// 		dryRun: true
// 	}
// 	unraid: {
// 		condition: {
// 			numDisks: 24,
// 			numProtected: 0,
// 		},
// 		disks: [
// 			{id: 1, name: "disk1", path: "/mnt/disk1"},
// 			{id: 2, name: "disk2", path: "/mnt/disk2"},
// 			{id: 3, name: "disk3", path: "/mnt/disk3"},
// 		],
// 		bytesToMove: 0,
// 		inProgress: false, // need to review this variable
// 	}
// 	opInProgress: null,
//  consoleLines: []
// }


// We define mapStateToProps where we'd normally use
// the @connect decorator so the data requirements are clear upfront, but then
// export the decorated component after the main class definition so
// the component can be tested w/ and w/o being connected.
// See: http://rackt.github.io/redux/docs/recipes/WritingTests.html
const mapStateToProps = (state) => ({
	unraid: state.unraid,
	config: state.config,
	opInProgress: state.opInProgress,
	consoleLines: state.consoleLines,
})
export class Home extends Component {
	static propTypes = {
		unraid: PropTypes.object.isRequired,
		config: PropTypes.object.isRequired,
		opInProgress: PropTypes.string.isRequired,
		consoleLines: PropTypes.array.isRequired,
		getConfig: PropTypes.func.isRequired,
		calculate: PropTypes.func.isRequired,
		move: PropTypes.func.isRequired,
	}

	render () {
		const { dispatch, unraid, config, opInProgress, consoleLines, getConfig, calculate, move } = this.props

		const ok = unraid.conditon.state === "STARTED"

		var warning = null;
		if (!ok) {
			return (
				<Panel msg="The array is not operational. Please start the array first." />
			)
		}

		var console = null;
		if (consoleLines.length !== 0) {
			return (
				<Console lines={consoleLines} />
			)
		}

		return (
			<section className="row">
				{ warning }

				<DashboardHeader ok={ok} {...this.props} />

				<Console />

				<DashboardContent ok={ok} {...this.props} />
			</section>
		)
	}
}

export default connect(mapStateToProps, actions)(Home)
