import React, { PureComponent } from 'react'
import { PropTypes } from 'prop-types'

import TreeMenu from 'react-tree-menu'
import classNames from 'classnames/bind'

import Wizard from './wizard'
import styles from '../styles/core.scss'

require('./tree-view.css')

const cx = classNames.bind(styles)

export default class GatherSource extends PureComponent {
	static propTypes = {
		store: PropTypes.object.isRequired,
		match: PropTypes.object.isRequired,
	}

	componentDidMount() {
		const {
			store: { state, actions },
		} = this.props

		actions.clearGatherConsole()
		if (Object.keys(state.gather.chosen).length === 0) {
			actions.getEntries()
		}
	}

	onCollapse = node => {
		const { gatherTreeCollapsed } = this.props.store.actions
		gatherTreeCollapsed(node)
	}

	onCheck = node => {
		const { gatherTreeChecked } = this.props.store.actions
		gatherTreeChecked(node)
	}

	render() {
		const {
			match,
			store: { state },
		} = this.props

		const header = (
			<section className={cx('row', 'bottom-spacer-half')}>
				<div className={cx('col-xs-12')}>
					<div className={cx('row')}>
						<div className={cx('col-xs-12')}>
							<div className={cx('flexSection', 'gatherSource')}>
								<div className={cx('col-xs-6')}>
									<span>SHARES EXPLORER</span>
								</div>
								<div className={cx('col-xs-6')}>
									<span>SELECTED</span>
								</div>
							</div>
						</div>
					</div>
				</div>
			</section>
		)

		const grid = (
			<section className={cx('row', 'bottom-spacer-half')}>
				<div className={cx('col-xs-12')}>
					<div className={cx('row')}>
						<div className={cx('col-xs-6', 'scroller')}>
							<TreeMenu
								expandIconClass="fa fa-chevron-right"
								collapseIconClass="fa fa-chevron-down"
								onTreeNodeCollapseChange={this.onCollapse}
								onTreeNodeCheckChange={this.onCheck}
								collapsible
								collapsed={false}
								data={state.gather.items}
							/>
						</div>

						<div className={cx('col-xs-6', 'scroller')}>
							<ul>
								{Object.keys(state.gather.chosen).map(chosen => (
									<li key={chosen}>
										- {chosen.slice(10)}:{' '}
										{state.gather.location && state.gather.location.presence[chosen.slice(10)]}
									</li>
								))}
							</ul>
						</div>
					</div>
				</div>
			</section>
		)

		return (
			<div>
				<Wizard match={match} store={this.props.store} />
				{header}
				{grid}
			</div>
		)
	}
}
