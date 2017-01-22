import React, { Component } from 'react'
import { Link } from 'react-router'
import 'font-awesome-webpack'

import styles from '../styles/core.scss'
import classNames from 'classnames/bind'

let cx = classNames.bind(styles)

export default class Settings extends Component {
	constructor(props) {
		super(props)

		this.state = {
			reservedAmount: props.store.state.config.reservedAmount,
			reservedUnit: props.store.state.config.reservedUnit,
			rsyncFlags: props.store.state.config.rsyncFlags,
		}
	}
	// componentDidMount() {
	// 	let { actions, dispatch } = this.props.store
	// 	dispatch(actions.getConfig)
	// }

						// <p>Set up the minimum amount of space that should be left free on each disk, after moving folders.</p>


	componentWillReceiveProps(next) {
		const { reservedAmount, reservedUnit, rsyncFlags } = next.store.state.config
		if (reservedAmount !== this.state.reservedAmount || reservedUnit !== this.state.reservedUnit ||  rsyncFlags !== this.state.rsyncFlags) {
			this.setState({
				reservedUnit,
				reservedAmount,
				rsyncFlags,
			})
		}
	}

	render() {
		// let { dispatch, state } = this.props
		// console.log('settings.render: ', this.props.store)
		let { state, actions } = this.props.store

		if (!state.config) {
			return null
		}

		// console.log('state.unraid: ', state.unraid)
		if (!state.unraid) {
			// dispatch(actions.getStorage)
			return null
		}

		const stateOk = state.unraid && state.unraid.condition.state === "STARTED"
		if (!stateOk) {
			// console.log('stateOk: ', stateOk)
			return (
				<section className={cx('row', 'bottom-spacer-half')}>
					<div className={cx('col-xs-12')}>
						<p className={cx('bg-warning')}>&nbsp; The array is not started. Please start the array before perfoming any operations with unBALANCE.</p>
					</div>
				</section>
			)
		}

		if (state.opInProgress === actions.calculate || state.opInProgress === actions.move) {
			return (
				<section className={cx('row', 'bottom-spacer-half')}>
					<div className={cx('col-xs-12')}>
						<p className={cx('bg-warning')}>&nbsp; {state.opInProgress} operation is currently under way. Wait until the operation has finished to make any settings changes.</p>
					</div>
				</section>
			)
		}

		let flags = this.state.rsyncFlags.join(' ')

		return (
			<div>

			<section className={cx('row', 'bottom-spacer-2x')}>
				<div className={cx('col-xs-12')}>
					<div>
						<h3>SET UP NOTIFICATIONS</h3>

						<p>Notifications rely on unRAID's notifications settings, so you need to set up unRAID first, in order to receive notifications from unBALANCE.</p>

						<span> Calculate: </span>
						<input id="calc0" className={cx('lspacer')} type="radio" name="calc" checked={state.config.notifyCalc === 0} onChange={this._setNotifyCalc(0)} />
						<label id="calc0" >No Notifications</label>

						<input id="calc1" className={cx('lspacer')} type="radio" name="calc" checked={state.config.notifyCalc === 1} onChange={this._setNotifyCalc(1)} />
						<label id="calc1" >Basic</label>

						<input id="calc2" className={cx('lspacer')} type="radio" name="calc" checked={state.config.notifyCalc === 2} onChange={this._setNotifyCalc(2)} />
						<label id="calc2" >Detailed</label>

						<br />

						<span> Move: </span>
						<input id="move0" className={cx('lspacer')} type="radio" name="move" checked={state.config.notifyMove === 0} onChange={this._setNotifyMove(0)} />
						<label id="move0">No Notifications</label>

						<input id="move0" className={cx('lspacer')} type="radio" name="move" checked={state.config.notifyMove === 1} onChange={this._setNotifyMove(1)} />
						<label id="move0">Basic</label>

						<input id="move0" className={cx('lspacer')} type="radio" name="move" checked={state.config.notifyMove === 2} onChange={this._setNotifyMove(2)} />
						<label id="move0">Detailed</label>
					</div>
				</div>
			</section>

			<section className={cx('row', 'bottom-spacer-2x')}>
				<div className={cx('col-xs-12')}>
					<div>
						<h3>RESERVED SPACE</h3>

						<p>unBALANCE uses the threshold defined here as the minimum free space that should be kept available in a target disk, when calculating how much the disk can be filled.</p>
						<p>This threshold cannot be less than 450Mb (hard limit set by this app).</p>

						<div className={cx('row')}>
							<div className={cx('col-xs-2')}>
								<div className={cx('addon')}>
									<input className={cx('addon-field')} type="number" value={this.state.reservedAmount} onChange={this._setReservedAmount} />
									<select className={cx('addon-item')} name="unit" value={this.state.reservedUnit} onChange={this._setReservedUnit}>
										<option value="%">%</option>
										<option value="Mb">Mb</option>
										<option value="Gb">Gb</option>
									</select>
								</div>
							</div>
							<div className={cx('col-xs-1')}>
								<button className={cx('btn', 'btn-primary')} onClick={this._setReservedSpace}>Apply</button>
							</div>
						</div>
					</div>
				</div>
			</section>

			<section className={cx('row', 'bottom-spacer-2x')}>
				<div className={cx('col-xs-12')}>
					<div>
						<h3>CUSTOM RSYNC FLAGS</h3>

						<p>Internally unBALANCE uses rsync to transfer files across disks.</p>
						<p>By default, rsync is invoked with <b>-avRX --partial</b> flags.</p>
						<p>Here you can set custom flags to override the default ones, except for the dry run flag which will be automatically added, if needed.</p>
						<p>It's strongly recommended to keep the -R flag, for optimal operation.</p>
						<p>Be careful with the flags you choose, since it can drastically alter the expected behaviour of rsync under unBALANCE.</p>

						<div className={cx('row')}>
							<div className={cx('col-xs-2')}>
								<div className={cx('addon')}>
									<input className={cx('addon-field')} type="string" value={flags} onChange={this._onChangeRsyncFlags} />
								</div>
							</div>
							<div className={cx('col-xs-4')}>
								<button className={cx('btn', 'btn-primary')} onClick={this._setRsyncFlags}>Apply</button>
								&nbsp;
								<button className={cx('btn', 'btn-primary')} onClick={this._setRsyncDefault}>Reset to default</button>
							</div>
						</div>
					</div>
				</div>
			</section>

			</div>
		)
	}

	// _addFolder(dispatch, e) {
	// 	console.log('key - value: ', e.key, e.target.value)
	// 	if (e.key !== "Enter") {
	// 		return
	// 	}

	// 	e.preventDefault()

	// 	dispatch(C.ADD_FOLDER, e.target.value)
	// }

	_setNotifyCalc = (notify) => (e) => {
		const { setNotifyCalc } = this.props.store.actions
		setNotifyCalc(notify)
	}

	_setNotifyMove = (notify) => (e) => {
		const { setNotifyMove } = this.props.store.actions
		setNotifyMove(notify)
	}

	_setReservedAmount = (e) => {
		this.setState({
			reservedAmount: e.target.value
		})
	}

	_setReservedUnit = (e) => {
		this.setState({
			reservedUnit: e.target.value
		})
	}

	_setReservedSpace = (e) => {
		const { setReservedSpace } = this.props.store.actions
		setReservedSpace(this.state.reservedAmount, this.state.reservedUnit)
	}

	_onChangeRsyncFlags = (e) => {
		this.setState({
			rsyncFlags: e.target.value.split(' ')
		})
	}

	_setRsyncFlags = () => {
		const { setRsyncFlags } = this.props.store.actions
		const flags = this.state.rsyncFlags.join(' ')
		setRsyncFlags(flags.trim().split(' '))
	}

	_setRsyncDefault = (e) => {
		const { setRsyncFlags } = this.props.store.actions
		setRsyncFlags(['-avPRX'])
	}
}
