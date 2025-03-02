<template>
	<div>
		<button
			class="btn btn-link pe-0 text-decoration-none link-dark text-nowrap"
			data-bs-toggle="modal"
			data-bs-target="#savingsModal"
		>
			<span class="d-inline d-sm-none">{{
				$t("footer.savings.footerShort", { percent })
			}}</span
			><span class="d-none d-sm-inline">{{
				$t("footer.savings.footerLong", { percent })
			}}</span
			><fa-icon icon="sun" class="icon ms-2 text-evcc"></fa-icon>
		</button>
		<div
			id="savingsModal"
			ref="modal"
			class="modal fade"
			tabindex="-1"
			role="dialog"
			aria-hidden="true"
		>
			<div class="modal-dialog modal-dialog-centered modal-dialog-scrollable" role="document">
				<div class="modal-content">
					<div class="modal-header">
						<h5 class="modal-title">
							<span class="d-block d-sm-none">
								{{
									$t("footer.savings.modalTitleShort", {
										percent,
										total: fmtKw(totalCharged * 1000, true, false),
									})
								}}
							</span>
							<span class="d-none d-sm-block">
								{{
									$t("footer.savings.modalTitleLong", {
										percent,
										total: fmtKw(totalCharged * 1000, true, false),
									})
								}}
							</span>
						</h5>
						<button
							type="button"
							class="btn-close"
							data-bs-dismiss="modal"
							aria-label="Close"
						></button>
					</div>
					<div class="modal-body py-4">
						<div class="chart-container mb-3">
							<div class="chart-legend d-flex flex-wrap justify-content-between mb-1">
								<div class="text-nowrap">
									<fa-icon icon="square" class="text-evcc"></fa-icon>
									{{
										$t("footer.savings.modalChartSelf", {
											self: fmtKw(selfConsumptionCharged * 1000, true, false),
										})
									}}
								</div>
								<div class="text-nowrap">
									<fa-icon icon="square" class="text-grid"></fa-icon>
									{{
										$t("footer.savings.modalChartGrid", {
											grid: fmtKw(gridCharged * 1000, true, false),
										})
									}}
								</div>
							</div>
							<div
								class="chart d-flex justify-content-stretch mb-1 rounded overflow-hidden"
							>
								<div
									v-if="totalCharged > 0"
									class="chart-item chart-item--self d-flex justify-content-center text-white flex-shrink-1"
									:style="{ width: `${percent}%` }"
								>
									<span class="text-truncate"> {{ percent }}% </span>
								</div>
								<div
									v-if="totalCharged > 0"
									class="chart-item chart-item--grid d-flex justify-content-center text-white flex-shrink-1"
									:style="{ width: `${100 - percent}%` }"
								>
									<span class="text-truncate"> {{ 100 - percent }}% </span>
								</div>
								<div
									v-if="totalCharged === 0"
									class="chart-item chart-item--no-data d-flex justify-content-center text-white w-100"
								>
									<span>{{ $t("footer.savings.modalNoData") }}</span>
								</div>
							</div>
						</div>
						<p class="mb-3">
							{{ $t("footer.savings.modalSavingsPrice") }}:
							<strong>{{ fmtPricePerKWh(effectivePrice, currency) }}</strong>
							<br />
							{{ $t("footer.savings.modalSavingsTotal") }}:
							<strong>{{ fmtMoney(amount, currency) }}</strong>
						</p>

						<p class="small text-muted mb-3">
							<a
								href="https://docs.evcc.io/docs/guides/setup/#ersparnisberechnung"
								target="_blank"
								class="text-muted"
							>
								{{ $t("footer.savings.modalExplaination") }}</a
							>:
							<span class="text-nowrap">
								{{
									$t("footer.savings.modalExplainationGrid", {
										gridPrice: fmtPricePerKWh(gridPrice, currency),
									})
								}}</span
							>,
							<span class="text-nowrap">
								{{
									$t("footer.savings.modalExplainationFeedIn", {
										feedInPrice: fmtPricePerKWh(feedInPrice, currency),
									})
								}}
							</span>
							<br />
							{{
								$t("footer.savings.modalServerStart", {
									since: fmtTimeAgo(secondsSinceStart()),
								})
							}}
						</p>

						<hr class="mb-4" />

						<Sponsor :sponsor="sponsor" class="mb-4" />

						<p class="small text-muted mb-0">
							<strong class="text-primary">
								<fa-icon icon="flask"></fa-icon>
								{{ $t("footer.savings.experimentalLabel") }}:
							</strong>
							{{ $t("footer.savings.experimentalText") }}
							<a
								href="https://github.com/evcc-io/evcc/discussions/2104"
								target="_blank"
								>GitHub Discussions</a
							>.
						</p>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<script>
import formatter from "../mixins/formatter";
import Sponsor from "./Sponsor.vue";

export default {
	name: "Savings",
	components: { Sponsor },
	mixins: [formatter],
	props: {
		selfConsumptionPercent: Number,
		since: { type: Number, default: 0 },
		sponsor: String,
		amount: { type: Number, default: 0 },
		effectivePrice: { type: Number, default: 0 },
		totalCharged: { type: Number, default: 0 },
		gridCharged: { type: Number, default: 0 },
		selfConsumptionCharged: { type: Number, default: 0 },
		gridPrice: { type: Number },
		feedInPrice: { type: Number },
		currency: String,
	},
	computed: {
		percent() {
			return Math.round(this.selfConsumptionPercent) || 0;
		},
	},
	methods: {
		secondsSinceStart() {
			return this.since * 1000 - Date.now();
		},
	},
};
</script>
<style scoped>
/* make modal a bottom drawer on small screens */
@media (max-width: 575px) {
	.modal-dialog.modal-dialog-centered {
		align-items: flex-end;
		margin-bottom: 0;
	}
	.modal.fade .modal-dialog {
		transition: transform 0.4s ease;
		transform: translate(0, 150px);
	}
	.modal.show .modal-dialog {
		transform: none;
	}
	.modal-dialog-scrollable {
		height: calc(100% - 0.5rem);
	}
	.modal-content {
		border-radius: 1rem 1rem 0 0;
	}
}

.chart {
	height: 1.6rem;
}

.chart-item--self {
	background-color: var(--evcc-self);
}
.chart-item--grid {
	background-color: var(--evcc-grid);
}
.chart-item--no-data {
	background-color: var(--bs-gray);
}

.chart-item {
	transition-property: width;
	transition-duration: 500ms;
	transition-timing-function: linear;
}
</style>
