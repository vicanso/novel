import request from 'axios';
import _ from 'lodash';

import {
  SERIES,
  INFLUX_RECORDS,
} from '../../urls';
import {
  getDate,
} from '../../helpers/util';

export default {
  data() {
    return {
      form: {},
      categories: null,
      trackers: null,
    };
  },
  methods: {
    async search() {
      const {
        category,
        account,
        begin,
        end,
      } = this.form;
      const url = INFLUX_RECORDS.replace(':measurement', 'userTracker');
      if (!category && !account) {
        this.$error('category and account can\'t be null');
        return;
      }
      const close = this.$loading();
      this.trackers = null;
      try {
        const conditions = {};
        if (category) {
          conditions.category = category;
        }
        if (account) {
          conditions.account = account;
        }
        const params = {
          'cache-control': 'no-cache',
        };
        if (!_.isEmpty(conditions)) {
          params.conditions = JSON.stringify(conditions);
        }
        if (begin) {
          params.start = new Date(begin).toISOString();
        }
        if (end) {
          params.end = new Date(end).toISOString();
        }
        const res = await request.get(url, {
          params,
        });
        this.trackers = _.map(res.data.trackers, (item) => {
          const params = JSON.parse(item.params);
          item.params = JSON.stringify(params, null, 2);
          item.date = getDate(item.time);
          return item;
        })
      } catch (err) {
        this.$error(err);
      } finally {
        close();
      }
    },
  },
  async beforeMount() {
    const url = SERIES.replace(':measurement', 'userTracker');
    const close = this.$loading();
    try {
      const res = await request.get(url);
      this.categories = res.data.category;
    } catch (err) {
      this.$error(err);
    } finally {
      close();
    }
  },
};
