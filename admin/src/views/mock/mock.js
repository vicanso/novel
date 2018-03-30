import request from 'axios';
import _ from 'lodash';

import {
  MOCKS,
} from '../../urls';
import {
  getDate,
  diff,
} from '../../helpers/util';

export default {
  data() {
    return {
      form: null,
      mocks: null,
      mode: -1,
    };
  },
  methods: {
    async loadingMocks() {
      const close = this.$loading();
      try {
        const res = await request.get(MOCKS);
        this.mocks = _.map(res.data.list, (item) => {
          // eslint-disable-next-line
          item.response = JSON.stringify(item.response, null, 2);
          item.date = getDate(item.updatedAt).substring(0, 16);
          return item;
        });
        this.mode = 0;
      } catch (err) {
        this.$error(err);
      } finally {
        close();
      }
    },
    edit(id) {
      this.form = _.clone(_.find(this.mocks, item => item.id === id));
      this.mode = 1;
    },
    add() {
      this.form = {};
      this.mode = 1;
    },
    async submit() {
      const {
        form,
      } = this;
      const keys = [
        'account',
        'url',
        'status',
        'response',
        'disabled',
        'description',
      ];
      const close = this.$loading();
      try {
        if (form.id) {
          const {
            id,
          } = form;
          const found = _.find(this.mocks, item => item.id === id);
          const updateData = diff(found, form, keys);
          await request.patch(`${MOCKS}/${id}`, updateData);
        } else {
          const data = {};
          _.forEach(keys, (key) => {
            data[key] = form[key];
          });
          await request.post(MOCKS, data);
        }
        await this.loadingMocks();
      } catch (err) {
        this.$error(err);
      } finally {
        close();
      }
    },
  },
  beforeMount() {
    this.loadingMocks();
  },
};
