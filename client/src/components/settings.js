import React, { PureComponent } from 'react'
import { PropTypes } from 'prop-types'

import 'font-awesome-webpack'
import classNames from 'classnames/bind'

import styles from '../styles/core.scss'

const cx = classNames.bind(styles)

export default class Settings extends PureComponent {
	static propTypes = {
		store: PropTypes.arrayOf(PropTypes.any).isRequired,
		actions: PropTypes.objectOf(PropTypes.func).isRequired,
	}

	constructor(props) {
		super(props)

		this.state = {
			reservedAmount: props.store.state.config.reservedAmount,
			reservedUnit: props.store.state.config.reservedUnit,
			rsyncFlags: props.store.state.config.rsyncFlags,
		}
	}

	componentWillReceiveProps(next) {
		const { reservedAmount, reservedUnit, rsyncFlags } = next.store.state.config
		if (
			reservedAmount !== this.state.reservedAmount ||
			reservedUnit !== this.state.reservedUnit ||
			rsyncFlags !== this.state.rsyncFlags
		) {
			this.setState({
				reservedUnit,
				reservedAmount,
				rsyncFlags,
			})
		}
	}

	setNotifyCalc = notify => () => {
		const { setNotifyCalc } = this.props.store.actions
		setNotifyCalc(notify)
	}

	setNotifyMove = notify => () => {
		const { setNotifyMove } = this.props.store.actions
		setNotifyMove(notify)
	}

	setReservedAmount = e => {
		this.setState({
			reservedAmount: e.target.value,
		})
	}

	setReservedUnit = e => {
		this.setState({
			reservedUnit: e.target.value,
		})
	}

	setReservedSpace = () => {
		const { setReservedSpace } = this.props.store.actions
		setReservedSpace(this.state.reservedAmount, this.state.reservedUnit)
	}

	_onChangeRsyncFlags = e => {
		this.setState({
			rsyncFlags: e.target.value.split(' '),
		})
	}

	_setRsyncFlags = () => {
		const { setRsyncFlags } = this.props.store.actions
		const flags = this.state.rsyncFlags.join(' ')
		setRsyncFlags(flags.trim().split(' '))
	}

	_setRsyncDefault = () => {
		const { setRsyncFlags } = this.props.store.actions
		setRsyncFlags(['-avPRX'])
	}

	setVerbosity = verbosity => () => {
		const { setVerbosity } = this.props.store.actions
		setVerbosity(verbosity)
	}

	render() {
		const { state, actions } = this.props.store

		if (!state.config) {
			return null
		}

		if (!state.unraid) {
			return null
		}

		const stateOk = state.unraid && state.unraid.condition.state === 'STARTED'
		if (!stateOk) {
			return (
				<section className={cx('row', 'bottom-spacer-half')}>
					<div className={cx('col-xs-12')}>
						<p className={cx('bg-warning')}>
							&nbsp; The array is not started. Please start the array before perfoming any operations with
							unBALANCE.
						</p>
					</div>
				</section>
			)
		}

		if (state.opInProgress === actions.calculate || state.opInProgress === actions.move) {
			return (
				<section className={cx('row', 'bottom-spacer-half')}>
					<div className={cx('col-xs-12')}>
						<p className={cx('bg-warning')}>
							&nbsp; {state.opInProgress} operation is currently under way. Wait until the operation has
							finished to make any settings changes.
						</p>
					</div>
				</section>
			)
		}

		const flags = this.state.rsyncFlags.join(' ')

		return (
			<div>
				<section className={cx('row', 'bottom-spacer-2x')}>
					<div className={cx('col-xs-12')}>
						<div>
							<h3>SET UP NOTIFICATIONS</h3>

							<p>
								Notifications rely on unRAID&apos;s notifications settings, so you need to set up unRAID
								first, in order to receive notifications from unBALANCE.
							</p>

							<span> Calculate: </span>
							<input
								id="calc0"
								className={cx('lspacer')}
								type="radio"
								name="calc"
								checked={state.config.notifyCalc === 0}
								onChange={this.setNotifyCalc(0)}
							/>
							<label htmlFor="calc0">No Notifications</label>

							<input
								id="calc1"
								className={cx('lspacer')}
								type="radio"
								name="calc"
								checked={state.config.notifyCalc === 1}
								onChange={this.setNotifyCalc(1)}
							/>
							<label htmlFor="calc1">Basic</label>

							<input
								id="calc2"
								className={cx('lspacer')}
								type="radio"
								name="calc"
								checked={state.config.notifyCalc === 2}
								onChange={this.setNotifyCalc(2)}
							/>
							<label htmlFor="calc2">Detailed</label>

							<br />

							<span> Transfer: </span>
							<input
								id="move0"
								className={cx('lspacer')}
								type="radio"
								name="move"
								checked={state.config.notifyMove === 0}
								onChange={this.setNotifyMove(0)}
							/>
							<label htmlFor="move0">No Notifications</label>

							<input
								id="move0"
								className={cx('lspacer')}
								type="radio"
								name="move"
								checked={state.config.notifyMove === 1}
								onChange={this.setNotifyMove(1)}
							/>
							<label htmlFor="move0">Basic</label>

							<input
								id="move0"
								className={cx('lspacer')}
								type="radio"
								name="move"
								checked={state.config.notifyMove === 2}
								onChange={this.setNotifyMove(2)}
							/>
							<label htmlFor="move0">Detailed</label>
						</div>
					</div>
				</section>

				<section className={cx('row', 'bottom-spacer-2x')}>
					<div className={cx('col-xs-12')}>
						<div>
							<h3>RESERVED SPACE</h3>

							<p>
								unBALANCE uses the threshold defined here as the minimum free space that should be kept
								available in a target disk, when calculating how much the disk can be filled.
							</p>
							<p>This threshold cannot be less than 450Mb (hard limit set by this app).</p>

							<div className={cx('row')}>
								<div className={cx('col-xs-2')}>
									<div className={cx('addon')}>
										<input
											className={cx('addon-field')}
											type="number"
											value={this.state.reservedAmount}
											onChange={this.setReservedAmount}
										/>
										<select
											className={cx('addon-item')}
											name="unit"
											value={this.state.reservedUnit}
											onChange={this.setReservedUnit}
										>
											<option value="%">%</option>
											<option value="Mb">Mb</option>
											<option value="Gb">Gb</option>
										</select>
									</div>
								</div>
								<div className={cx('col-xs-1')}>
									<button className={cx('btn', 'btn-primary')} onClick={this.setReservedSpace}>
										Apply
									</button>
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
							<p>
								By default, rsync is invoked with <b>-avPRX</b> flags.
							</p>
							<p>
								Here you can set custom flags to override the default ones, except for the dry run flag
								which will be automatically added, if needed.
							</p>
							<p>It&apos;s strongly recommended to keep the -R flag, for optimal operation.</p>
							<p>
								Be careful with the flags you choose, since it can drastically alter the expected
								behaviour of rsync under unBALANCE.
							</p>
							<p>
								<span className={cx('opWarning')}>
									Also note that for proper VALIDATE functionality, the custom flags MUST being with
									&quot;-a&quot;.
								</span>
							</p>

							<div className={cx('row')}>
								<div className={cx('col-xs-2')}>
									<div className={cx('addon')}>
										<input
											className={cx('addon-field')}
											type="string"
											value={flags}
											onChange={this.onChangeRsyncFlags}
										/>
									</div>
								</div>
								<div className={cx('col-xs-4')}>
									<button className={cx('btn', 'btn-primary')} onClick={this.setRsyncFlags}>
										Apply
									</button>
									&nbsp;
									<button className={cx('btn', 'btn-primary')} onClick={this.setRsyncDefault}>
										Reset to default
									</button>
								</div>
							</div>
						</div>
					</div>
				</section>

				<section className={cx('row', 'bottom-spacer-2x')}>
					<div className={cx('col-xs-12')}>
						<div>
							<h3>SET LOG VERBOSITY</h3>

							<p>
								Full verbosity will print each line generated in the transfer (rsync) phase. Normal
								verbosity will not, thus greatly reducing the amount of logging.
							</p>

							<span> Verbosity: </span>
							<input
								id="verb0"
								className={cx('lspacer')}
								type="radio"
								name="verb"
								checked={state.config.verbosity === 0}
								onChange={this.setVerbosity(0)}
							/>
							<label htmlFor="verb0">Normal</label>

							<input
								id="verb1"
								className={cx('lspacer')}
								type="radio"
								name="verb"
								checked={state.config.verbosity === 1}
								onChange={this.setVerbosity(1)}
							/>
							<label htmlFor="verb1">Full</label>
						</div>
					</div>
				</section>
			</div>
		)
	}
}
