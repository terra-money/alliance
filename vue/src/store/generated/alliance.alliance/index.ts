import { Client, registry, MissingWalletError } from 'alliance-client-ts'

import { AllianceAsset } from "alliance-client-ts/alliance.alliance/types"
import { AddAssetProposal } from "alliance-client-ts/alliance.alliance/types"
import { RemoveAssetProposal } from "alliance-client-ts/alliance.alliance/types"
import { UpdateAssetProposal } from "alliance-client-ts/alliance.alliance/types"
import { QueuedRewardRateChange } from "alliance-client-ts/alliance.alliance/types"
import { RewardRateChangeSnapshot } from "alliance-client-ts/alliance.alliance/types"
import { Delegation } from "alliance-client-ts/alliance.alliance/types"
import { DelegationResponse } from "alliance-client-ts/alliance.alliance/types"
import { Redelegation } from "alliance-client-ts/alliance.alliance/types"
import { QueuedRedelegation } from "alliance-client-ts/alliance.alliance/types"
import { Undelegation } from "alliance-client-ts/alliance.alliance/types"
import { QueuedUndelegation } from "alliance-client-ts/alliance.alliance/types"
import { AllianceValidatorInfo } from "alliance-client-ts/alliance.alliance/types"
import { Params } from "alliance-client-ts/alliance.alliance/types"
import { RewardHistory } from "alliance-client-ts/alliance.alliance/types"
import { NewAllianceAssetMsg } from "alliance-client-ts/alliance.alliance/types"


export { AllianceAsset, AddAssetProposal, RemoveAssetProposal, UpdateAssetProposal, QueuedRewardRateChange, RewardRateChangeSnapshot, Delegation, DelegationResponse, Redelegation, QueuedRedelegation, Undelegation, QueuedUndelegation, AllianceValidatorInfo, Params, RewardHistory, NewAllianceAssetMsg };

function initClient(vuexGetters) {
	return new Client(vuexGetters['common/env/getEnv'], vuexGetters['common/wallet/signer'])
}

function mergeResults(value, next_values) {
	for (let prop of Object.keys(next_values)) {
		if (Array.isArray(next_values[prop])) {
			value[prop]=[...value[prop], ...next_values[prop]]
		}else{
			value[prop]=next_values[prop]
		}
	}
	return value
}

type Field = {
	name: string;
	type: unknown;
}
function getStructure(template) {
	let structure: {fields: Field[]} = { fields: [] }
	for (const [key, value] of Object.entries(template)) {
		let field = { name: key, type: typeof value }
		structure.fields.push(field)
	}
	return structure
}
const getDefaultState = () => {
	return {
				Params: {},
				Alliances: {},
				Alliance: {},
				AlliancesDelegation: {},
				AlliancesDelegationByValidator: {},
				AllianceDelegation: {},
				AllianceDelegationRewards: {},
				
				_Structure: {
						AllianceAsset: getStructure(AllianceAsset.fromPartial({})),
						AddAssetProposal: getStructure(AddAssetProposal.fromPartial({})),
						RemoveAssetProposal: getStructure(RemoveAssetProposal.fromPartial({})),
						UpdateAssetProposal: getStructure(UpdateAssetProposal.fromPartial({})),
						QueuedRewardRateChange: getStructure(QueuedRewardRateChange.fromPartial({})),
						RewardRateChangeSnapshot: getStructure(RewardRateChangeSnapshot.fromPartial({})),
						Delegation: getStructure(Delegation.fromPartial({})),
						DelegationResponse: getStructure(DelegationResponse.fromPartial({})),
						Redelegation: getStructure(Redelegation.fromPartial({})),
						QueuedRedelegation: getStructure(QueuedRedelegation.fromPartial({})),
						Undelegation: getStructure(Undelegation.fromPartial({})),
						QueuedUndelegation: getStructure(QueuedUndelegation.fromPartial({})),
						AllianceValidatorInfo: getStructure(AllianceValidatorInfo.fromPartial({})),
						Params: getStructure(Params.fromPartial({})),
						RewardHistory: getStructure(RewardHistory.fromPartial({})),
						NewAllianceAssetMsg: getStructure(NewAllianceAssetMsg.fromPartial({})),
						
		},
		_Registry: registry,
		_Subscriptions: new Set(),
	}
}

// initial state
const state = getDefaultState()

export default {
	namespaced: true,
	state,
	mutations: {
		RESET_STATE(state) {
			Object.assign(state, getDefaultState())
		},
		QUERY(state, { query, key, value }) {
			state[query][JSON.stringify(key)] = value
		},
		SUBSCRIBE(state, subscription) {
			state._Subscriptions.add(JSON.stringify(subscription))
		},
		UNSUBSCRIBE(state, subscription) {
			state._Subscriptions.delete(JSON.stringify(subscription))
		}
	},
	getters: {
				getParams: (state) => (params = { params: {}}) => {
					if (!(<any> params).query) {
						(<any> params).query=null
					}
			return state.Params[JSON.stringify(params)] ?? {}
		},
				getAlliances: (state) => (params = { params: {}}) => {
					if (!(<any> params).query) {
						(<any> params).query=null
					}
			return state.Alliances[JSON.stringify(params)] ?? {}
		},
				getAlliance: (state) => (params = { params: {}}) => {
					if (!(<any> params).query) {
						(<any> params).query=null
					}
			return state.Alliance[JSON.stringify(params)] ?? {}
		},
				getAlliancesDelegation: (state) => (params = { params: {}}) => {
					if (!(<any> params).query) {
						(<any> params).query=null
					}
			return state.AlliancesDelegation[JSON.stringify(params)] ?? {}
		},
				getAlliancesDelegationByValidator: (state) => (params = { params: {}}) => {
					if (!(<any> params).query) {
						(<any> params).query=null
					}
			return state.AlliancesDelegationByValidator[JSON.stringify(params)] ?? {}
		},
				getAllianceDelegation: (state) => (params = { params: {}}) => {
					if (!(<any> params).query) {
						(<any> params).query=null
					}
			return state.AllianceDelegation[JSON.stringify(params)] ?? {}
		},
				getAllianceDelegationRewards: (state) => (params = { params: {}}) => {
					if (!(<any> params).query) {
						(<any> params).query=null
					}
			return state.AllianceDelegationRewards[JSON.stringify(params)] ?? {}
		},
				
		getTypeStructure: (state) => (type) => {
			return state._Structure[type].fields
		},
		getRegistry: (state) => {
			return state._Registry
		}
	},
	actions: {
		init({ dispatch, rootGetters }) {
			console.log('Vuex module: alliance.alliance initialized!')
			if (rootGetters['common/env/client']) {
				rootGetters['common/env/client'].on('newblock', () => {
					dispatch('StoreUpdate')
				})
			}
		},
		resetState({ commit }) {
			commit('RESET_STATE')
		},
		unsubscribe({ commit }, subscription) {
			commit('UNSUBSCRIBE', subscription)
		},
		async StoreUpdate({ state, dispatch }) {
			state._Subscriptions.forEach(async (subscription) => {
				try {
					const sub=JSON.parse(subscription)
					await dispatch(sub.action, sub.payload)
				}catch(e) {
					throw new Error('Subscriptions: ' + e.message)
				}
			})
		},
		
		
		
		 		
		
		
		async QueryParams({ commit, rootGetters, getters }, { options: { subscribe, all} = { subscribe:false, all:false}, params, query=null }) {
			try {
				const key = params ?? {};
				const client = initClient(rootGetters);
				let value= (await client.AllianceAlliance.query.queryParams()).data
				
					
				commit('QUERY', { query: 'Params', key: { params: {...key}, query}, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryParams', payload: { options: { all }, params: {...key},query }})
				return getters['getParams']( { params: {...key}, query}) ?? {}
			} catch (e) {
				throw new Error('QueryClient:QueryParams API Node Unavailable. Could not perform query: ' + e.message)
				
			}
		},
		
		
		
		
		 		
		
		
		async QueryAlliances({ commit, rootGetters, getters }, { options: { subscribe, all} = { subscribe:false, all:false}, params, query=null }) {
			try {
				const key = params ?? {};
				const client = initClient(rootGetters);
				let value= (await client.AllianceAlliance.query.queryAlliances(query ?? undefined)).data
				
					
				while (all && (<any> value).pagination && (<any> value).pagination.next_key!=null) {
					let next_values=(await client.AllianceAlliance.query.queryAlliances({...query ?? {}, 'pagination.key':(<any> value).pagination.next_key} as any)).data
					value = mergeResults(value, next_values);
				}
				commit('QUERY', { query: 'Alliances', key: { params: {...key}, query}, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryAlliances', payload: { options: { all }, params: {...key},query }})
				return getters['getAlliances']( { params: {...key}, query}) ?? {}
			} catch (e) {
				throw new Error('QueryClient:QueryAlliances API Node Unavailable. Could not perform query: ' + e.message)
				
			}
		},
		
		
		
		
		 		
		
		
		async QueryAlliance({ commit, rootGetters, getters }, { options: { subscribe, all} = { subscribe:false, all:false}, params, query=null }) {
			try {
				const key = params ?? {};
				const client = initClient(rootGetters);
				let value= (await client.AllianceAlliance.query.queryAlliance( key.denom)).data
				
					
				commit('QUERY', { query: 'Alliance', key: { params: {...key}, query}, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryAlliance', payload: { options: { all }, params: {...key},query }})
				return getters['getAlliance']( { params: {...key}, query}) ?? {}
			} catch (e) {
				throw new Error('QueryClient:QueryAlliance API Node Unavailable. Could not perform query: ' + e.message)
				
			}
		},
		
		
		
		
		 		
		
		
		async QueryAlliancesDelegation({ commit, rootGetters, getters }, { options: { subscribe, all} = { subscribe:false, all:false}, params, query=null }) {
			try {
				const key = params ?? {};
				const client = initClient(rootGetters);
				let value= (await client.AllianceAlliance.query.queryAlliancesDelegation( key.delegator_addr, query ?? undefined)).data
				
					
				while (all && (<any> value).pagination && (<any> value).pagination.next_key!=null) {
					let next_values=(await client.AllianceAlliance.query.queryAlliancesDelegation( key.delegator_addr, {...query ?? {}, 'pagination.key':(<any> value).pagination.next_key} as any)).data
					value = mergeResults(value, next_values);
				}
				commit('QUERY', { query: 'AlliancesDelegation', key: { params: {...key}, query}, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryAlliancesDelegation', payload: { options: { all }, params: {...key},query }})
				return getters['getAlliancesDelegation']( { params: {...key}, query}) ?? {}
			} catch (e) {
				throw new Error('QueryClient:QueryAlliancesDelegation API Node Unavailable. Could not perform query: ' + e.message)
				
			}
		},
		
		
		
		
		 		
		
		
		async QueryAlliancesDelegationByValidator({ commit, rootGetters, getters }, { options: { subscribe, all} = { subscribe:false, all:false}, params, query=null }) {
			try {
				const key = params ?? {};
				const client = initClient(rootGetters);
				let value= (await client.AllianceAlliance.query.queryAlliancesDelegationByValidator( key.delegator_addr,  key.validator_addr, query ?? undefined)).data
				
					
				while (all && (<any> value).pagination && (<any> value).pagination.next_key!=null) {
					let next_values=(await client.AllianceAlliance.query.queryAlliancesDelegationByValidator( key.delegator_addr,  key.validator_addr, {...query ?? {}, 'pagination.key':(<any> value).pagination.next_key} as any)).data
					value = mergeResults(value, next_values);
				}
				commit('QUERY', { query: 'AlliancesDelegationByValidator', key: { params: {...key}, query}, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryAlliancesDelegationByValidator', payload: { options: { all }, params: {...key},query }})
				return getters['getAlliancesDelegationByValidator']( { params: {...key}, query}) ?? {}
			} catch (e) {
				throw new Error('QueryClient:QueryAlliancesDelegationByValidator API Node Unavailable. Could not perform query: ' + e.message)
				
			}
		},
		
		
		
		
		 		
		
		
		async QueryAllianceDelegation({ commit, rootGetters, getters }, { options: { subscribe, all} = { subscribe:false, all:false}, params, query=null }) {
			try {
				const key = params ?? {};
				const client = initClient(rootGetters);
				let value= (await client.AllianceAlliance.query.queryAllianceDelegation( key.delegator_addr,  key.validator_addr,  key.denom, query ?? undefined)).data
				
					
				while (all && (<any> value).pagination && (<any> value).pagination.next_key!=null) {
					let next_values=(await client.AllianceAlliance.query.queryAllianceDelegation( key.delegator_addr,  key.validator_addr,  key.denom, {...query ?? {}, 'pagination.key':(<any> value).pagination.next_key} as any)).data
					value = mergeResults(value, next_values);
				}
				commit('QUERY', { query: 'AllianceDelegation', key: { params: {...key}, query}, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryAllianceDelegation', payload: { options: { all }, params: {...key},query }})
				return getters['getAllianceDelegation']( { params: {...key}, query}) ?? {}
			} catch (e) {
				throw new Error('QueryClient:QueryAllianceDelegation API Node Unavailable. Could not perform query: ' + e.message)
				
			}
		},
		
		
		
		
		 		
		
		
		async QueryAllianceDelegationRewards({ commit, rootGetters, getters }, { options: { subscribe, all} = { subscribe:false, all:false}, params, query=null }) {
			try {
				const key = params ?? {};
				const client = initClient(rootGetters);
				let value= (await client.AllianceAlliance.query.queryAllianceDelegationRewards( key.delegator_addr,  key.validator_addr,  key.denom, query ?? undefined)).data
				
					
				while (all && (<any> value).pagination && (<any> value).pagination.next_key!=null) {
					let next_values=(await client.AllianceAlliance.query.queryAllianceDelegationRewards( key.delegator_addr,  key.validator_addr,  key.denom, {...query ?? {}, 'pagination.key':(<any> value).pagination.next_key} as any)).data
					value = mergeResults(value, next_values);
				}
				commit('QUERY', { query: 'AllianceDelegationRewards', key: { params: {...key}, query}, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryAllianceDelegationRewards', payload: { options: { all }, params: {...key},query }})
				return getters['getAllianceDelegationRewards']( { params: {...key}, query}) ?? {}
			} catch (e) {
				throw new Error('QueryClient:QueryAllianceDelegationRewards API Node Unavailable. Could not perform query: ' + e.message)
				
			}
		},
		
		
		
		
	}
}
