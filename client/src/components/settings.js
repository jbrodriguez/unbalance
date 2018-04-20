import React, { PureComponent } from 'react'
import { PropTypes } from 'prop-types'

// import 'font-awesome-webpack'
import classNames from 'classnames/bind'

import styles from '../styles/core.scss'

const cx = classNames.bind(styles)

export default class Settings extends PureComponent {
	static propTypes = {
		store: PropTypes.objectOf(PropTypes.any).isRequired,
	}

	constructor(props) {
		super(props)

		this.state = {
			reservedAmount: props.store.state.config.reservedAmount,
			reservedUnit: props.store.state.config.reservedUnit,
			rsyncArgs: props.store.state.config.rsyncArgs,
		}
	}

	componentDidMount() {
		const { actions } = this.props.store
		actions.getConfig()
	}

	componentWillReceiveProps(next) {
		const { reservedAmount, reservedUnit, rsyncArgs } = next.store.state.config
		if (
			reservedAmount !== this.state.reservedAmount ||
			reservedUnit !== this.state.reservedUnit ||
			rsyncArgs !== this.state.rsyncArgs
		) {
			this.setState({
				reservedUnit,
				reservedAmount,
				rsyncArgs,
			})
		}
	}

	setNotifyPlan = notify => () => {
		const { setNotifyPlan } = this.props.store.actions
		setNotifyPlan(notify)
	}

	setNotifyTransfer = notify => () => {
		const { setNotifyTransfer } = this.props.store.actions
		setNotifyTransfer(notify)
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

	_onChangeRsyncArgs = e => {
		this.setState({
			rsyncArgs: e.target.value.split(' '),
		})
	}

	_setRsyncArgs = () => {
		const { setRsyncArgs } = this.props.store.actions
		const args = this.state.rsyncArgs.join(' ')
		setRsyncArgs(args.trim().split(' '))
	}

	_setRsyncDefault = () => {
		const { setRsyncArgs } = this.props.store.actions
		setRsyncArgs(['-X'])
	}

	setVerbosity = verbosity => () => {
		const { setVerbosity } = this.props.store.actions
		setVerbosity(verbosity)
	}

	setUpdateCheck = updateCheck => () => {
		const { setUpdateCheck } = this.props.store.actions
		setUpdateCheck(updateCheck)
	}

	render() {
		const { state } = this.props.store

		const args = this.state.rsyncArgs.join(' ')

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

							<span> Planning: </span>
							<input
								id="plan0"
								className={cx('lspacer')}
								type="radio"
								name="plan"
								checked={state.config.notifyPlan === 0}
								onChange={this.setNotifyPlan(0)}
							/>
							<label htmlFor="plan0">No Notifications</label>

							<input
								id="plan1"
								className={cx('lspacer')}
								type="radio"
								name="plan"
								checked={state.config.notifyPlan === 1}
								onChange={this.setNotifyPlan(1)}
							/>
							<label htmlFor="plan1">Basic</label>

							<input
								id="plan2"
								className={cx('lspacer')}
								type="radio"
								name="plan"
								checked={state.config.notifyPlan === 2}
								onChange={this.setNotifyPlan(2)}
							/>
							<label htmlFor="plan2">Detailed</label>

							<br />

							<span> Transfer: </span>
							<input
								id="transfer0"
								className={cx('lspacer')}
								type="radio"
								name="transfer"
								checked={state.config.notifyTransfer === 0}
								onChange={this.setNotifyTransfer(0)}
							/>
							<label htmlFor="transfer0">No Notifications</label>

							<input
								id="transfer0"
								className={cx('lspacer')}
								type="radio"
								name="transfer"
								checked={state.config.notifyTransfer === 1}
								onChange={this.setNotifyTransfer(1)}
							/>
							<label htmlFor="transfer0">Basic</label>

							<input
								id="transfer0"
								className={cx('lspacer')}
								type="radio"
								name="transfer"
								checked={state.config.notifyTransfer === 2}
								onChange={this.setNotifyTransfer(2)}
							/>
							<label htmlFor="transfer0">Detailed</label>
						</div>
					</div>
				</section>

				<section className={cx('row', 'bottom-spacer-2x')}>
					<div className={cx('col-xs-12')}>
						<div>
							<h3>RESERVED SPACE</h3>

							<p>
								unBALANCE uses the threshold defined here as the minimum free space that should be kept
								available in a target disk, when planning how much the disk can be filled.
							</p>
							<p>This threshold cannot be less than 512Mb (hard limit set by this app).</p>

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

							<p>Internally, unBALANCE uses rsync to transfer files across disks.</p>
							<p>
								By default, rsync is invoked with <b>-avPRX</b> flags. Note that the <b>X</b> flag is
								customizable, so you can remove it if needed.
							</p>
							<p>
								You can add custom flags, except for the dry run flag which will be automatically added,
								if needed.
							</p>
							<p>
								Be careful with the flags you choose, since it can drastically alter the expected
								behaviour of rsync under unBALANCE.
							</p>

							<div className={cx('row')}>
								<div className={cx('col-xs-2')}>
									<div className={cx('addon')}>
										<input
											className={cx('addon-field')}
											type="string"
											value={args}
											onChange={this._onChangeRsyncArgs}
										/>
									</div>
								</div>
								<div className={cx('col-xs-4')}>
									<button className={cx('btn', 'btn-primary')} onClick={this._setRsyncArgs}>
										Apply
									</button>
									&nbsp;
									<button className={cx('btn', 'btn-primary')} onClick={this._setRsyncDefault}>
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
								Full verbosity will affect logging in two ways:<br />
								- It will print each line generated in the transfer (rsync) phase.<br />
								- It will print each line generated while checking for permission issues.<br />
								Normal verbosity will not, thus greatly reducing the amount of logging.
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

				<section className={cx('row', 'bottom-spacer-2x')}>
					<div className={cx('col-xs-12')}>
						<div>
							<h3>CHECK FOR UPDATES</h3>

							<p>
								On will check for a plugin update on start.<br />
								Off disables the check.<br />
							</p>

							<span> Check: </span>
							<input
								id="check0"
								className={cx('lspacer')}
								type="radio"
								name="check"
								checked={state.config.checkForUpdate === 1}
								onChange={this.setUpdateCheck(1)}
							/>
							<label htmlFor="check0">On</label>

							<input
								id="check1"
								className={cx('lspacer')}
								type="radio"
								name="check"
								checked={state.config.checkForUpdate === 0}
								onChange={this.setUpdateCheck(0)}
							/>
							<label htmlFor="check1">Off</label>
						</div>
					</div>
				</section>
			</div>
		)
	}
}
